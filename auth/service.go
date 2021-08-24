package auth

import (
	"context"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
	"github.com/pkg/errors"
)

type Service interface{}

type service struct {
	authRepo Repository
}

func NewService(authRepo Repository) Service {
	return &service{
		authRepo: authRepo,
	}
}

func (s *service) SignUp(ctx context.Context, req *SignUpReq) (SignUpResp, error) {
	user, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		return SignUpResp{}, app.NewRestError(app.ErrInternalServerError, errors.Wrap(err, "authRepo.SignUp.FindByEmail"))
	}

	user = models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err = s.authRepo.Create(ctx, &user)
	if err != nil {
		return SignUpResp{}, app.NewRestError(app.ErrInternalServerError, errors.Wrap(err, "authRepo.SignUp.Create"))
	}

	res := SignUpResp{
		User: User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Picture:   user.Picture.String,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return res, nil
}
