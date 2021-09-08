package auth

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	jwt.StandardClaims
	UserID int64
}

type RefreshClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"userID"`
}

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignInReq struct {
	Email    string `json:"email" validate:"required,gte=5,lte=255"`
	Password string `json:"password" validate:"required,gte=5,lte=50"`
}

func (r *SignInReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SignInGoogleReq struct {
	GoogleAccessToken string `json:"googleAccessToken"`
}

func (r *SignInGoogleReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SignInResp struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type SignUpReq struct {
	Name     string `json:"name" validate:"required,gte=5,lte=50"`
	Email    string `json:"email" validate:"required,gte=5,lte=255"`
	Password string `json:"password" validate:"required,gte=5,lte=50"`
}

func (r *SignUpReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SignUpResp struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type GetUserResp User

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResp struct {
	UserID int64 `json:"userID"`
	Token  Token `json:"token"`
}
