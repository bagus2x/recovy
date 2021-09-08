package auth

import (
	"context"
	"strings"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/config"
	"github.com/bagus2x/recovy/models"
	"github.com/golang-jwt/jwt"
)

type Service interface {
	SignIn(ctx context.Context, req *SignInReq) (SignInResp, error)
	SignUp(ctx context.Context, req *SignUpReq) (SignUpResp, error)
	GetUserByID(ctx context.Context, userID int64) (GetUserResp, error)
	SignOut(ctx context.Context, userID int64) error
	RefreshToken(ctx context.Context, req *RefreshTokenReq) (RefreshTokenResp, error)
	ExtractAccessToken(tokenStr string) (AccessClaims, error)
}

type service struct {
	authRepo      Repository
	authCacheRepo CacheRepository
	cfg           *config.Config
}

func NewService(authRepo Repository, authCacheRepo CacheRepository, cfg *config.Config) Service {
	return &service{
		authRepo:      authRepo,
		authCacheRepo: authCacheRepo,
		cfg:           cfg,
	}
}

func (s *service) SignIn(ctx context.Context, req *SignInReq) (SignInResp, error) {
	err := req.Validate()
	if err != nil {
		return SignInResp{}, err
	}

	user, err := s.authRepo.FindByEmail(ctx, req.Email)
	if app.ErrorCode(err) == app.ENotFound {
		return SignInResp{}, app.NewError(nil, app.ENotFound, "User not found")
	} else if err != nil {
		return SignInResp{}, err
	}

	isMatch := user.ComparePasswords(req.Password)
	if !isMatch {
		return SignInResp{}, app.NewError(nil, app.EBadRequest, "Password does not match")
	}

	accessToken, err := s.createAccessToken(user.ID)
	if err != nil {
		return SignInResp{}, err
	}

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		return SignInResp{}, err
	}

	err = s.authCacheRepo.SetRefreshToken(user.ID, refreshToken)
	if err != nil {
		return SignInResp{}, err
	}

	res := SignInResp{
		Token: Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: User{
			ID:        user.ID,
			Picture:   user.Picture,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return res, nil
}

func (s *service) SignUp(ctx context.Context, req *SignUpReq) (SignUpResp, error) {
	err := req.Validate()
	if err != nil {
		return SignUpResp{}, err
	}

	user, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil && app.ErrorCode(err) != app.ENotFound {
		return SignUpResp{}, err
	}

	if user.Email == req.Email {
		return SignUpResp{}, app.NewError(nil, app.Econflict, "Email already exist")
	}

	user = models.User{
		Name:      req.Name,
		Password:  req.Password,
		Email:     req.Email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = user.HashPassword()
	if err != nil {
		return SignUpResp{}, err
	}

	err = s.authRepo.Create(ctx, &user)
	if err != nil {
		return SignUpResp{}, err
	}

	accessToken, err := s.createAccessToken(user.ID)
	if err != nil {
		return SignUpResp{}, err
	}

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		return SignUpResp{}, err
	}

	err = s.authCacheRepo.SetRefreshToken(user.ID, refreshToken)
	if err != nil {
		return SignUpResp{}, err
	}

	res := SignUpResp{
		Token: Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: User{
			ID:        user.ID,
			Picture:   user.Picture,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return res, nil
}

func (s *service) GetUserByID(ctx context.Context, userID int64) (GetUserResp, error) {
	p, err := s.authRepo.FindByID(ctx, userID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetUserResp{}, app.NewError(err, app.ENotFound, "User does not exist")
	} else if err != nil {
		return GetUserResp{}, err
	}

	res := GetUserResp{
		ID:        userID,
		Picture:   p.Picture,
		Name:      p.Name,
		Email:     p.Email,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.CreatedAt,
	}

	return res, nil
}

func (s *service) RefreshToken(ctx context.Context, req *RefreshTokenReq) (RefreshTokenResp, error) {
	claims, err := s.extractRefreshToken(req.RefreshToken)
	if err != nil {
		return RefreshTokenResp{}, err
	}

	refreshToken, err := s.authCacheRepo.GetRefreshToken(claims.UserID)
	if err != nil {
		return RefreshTokenResp{}, app.NewError(err, app.EForbidden)
	}

	if refreshToken != req.RefreshToken {
		return RefreshTokenResp{}, app.NewError(nil, app.EInvalidRefreshToken)
	}

	accessToken, err := s.createAccessToken(claims.UserID)
	if err != nil {
		return RefreshTokenResp{}, err
	}

	refreshToken, err = s.createRefreshToken(claims.UserID)
	if err != nil {
		return RefreshTokenResp{}, err
	}

	err = s.authCacheRepo.SetRefreshToken(claims.UserID, refreshToken)
	if err != nil {
		return RefreshTokenResp{}, app.NewError(err, app.EInternal)
	}

	res := RefreshTokenResp{
		UserID: claims.UserID,
		Token: Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return res, nil
}

func (s *service) SignOut(ctx context.Context, userID int64) error {
	affected := s.authCacheRepo.DeleteRefreshToken(userID)
	if !affected {
		return app.NewError(nil, app.ENotFound, "User does not exist")
	}

	return nil
}

func (s *service) createAccessToken(userID int64) (string, error) {
	claims := AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.cfg.AccessTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.AccessTokenKey()))
}

func (s *service) ExtractAccessToken(tokenStr string) (AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app.NewError(nil, app.EInternal, "Unexpected signing method")
		}
		return []byte(s.cfg.AccessTokenKey()), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return AccessClaims{}, app.NewError(err, app.EAccessTokenExpired, "Access token has expired")
		}
		return AccessClaims{}, app.NewError(err, app.EInvalidAccessToken, "Invalid access token")
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid && claims != nil {
		return *claims, nil
	}

	return AccessClaims{}, app.NewError(err, app.EInvalidAccessToken, "Invalid access token")
}

func (s *service) createRefreshToken(userID int64) (string, error) {
	claims := AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.cfg.RefreshTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.RefreshTokenKey()))
}

func (s *service) extractRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app.NewError(nil, app.EInternal, "Unexpected signing method")
		}
		return []byte(s.cfg.RefreshTokenKey()), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, app.NewError(err, app.ERefreshTokenExpired, "Refresh token has expired")
		}
		return nil, app.NewError(err, app.EInvalidRefreshToken, "Invalid refresh token")
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, app.NewError(err, app.EInvalidRefreshToken, "Invalid refresh token")
}
