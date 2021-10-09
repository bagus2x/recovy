package webinar

import (
	"context"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Service interface {
	Create(ctx context.Context, req *CreateWebinarReq) (CreateWebinarResp, error)
	GetByID(ctx context.Context, webinarID int64) (GetWebinarResp, error)
	GetByCategory(ctx context.Context, category string) ([]GetWebinarResp, error)
	Get(ctx context.Context) ([]GetWebinarResp, error)
	Delete(ctx context.Context, webinarID, authorID int64) error
}

type service struct {
	webinarRepo Repository
}

func NewService(webinarRepo Repository) Service {
	return &service{
		webinarRepo: webinarRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreateWebinarReq) (CreateWebinarResp, error) {
	err := req.Validate()
	if err != nil {
		return CreateWebinarResp{}, err
	}

	webinar := models.Webinar{
		Author: models.User{
			ID: req.AuthorID,
		},
		Picture:     req.Picture,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		StartDate:   req.StartDate,
		LastDate:    req.LastDate,
		Time:        req.Time,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = s.webinarRepo.Create(ctx, &webinar)
	if err != nil {
		return CreateWebinarResp{}, err
	}

	res := CreateWebinarResp{
		ID:          webinar.ID,
		AuthorID:    webinar.Author.ID,
		Picture:     webinar.Picture,
		Title:       webinar.Title,
		Description: webinar.Description,
		Category:    webinar.Category,
		StartDate:   webinar.StartDate,
		LastDate:    webinar.LastDate,
		Time:        webinar.Time,
		CreatedAt:   webinar.CreatedAt,
		UpdatedAt:   webinar.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, webinarID int64) (GetWebinarResp, error) {
	webinar, err := s.webinarRepo.FindByID(ctx, webinarID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetWebinarResp{}, app.NewError(err, app.ENotFound, "Webinar not found")
	} else if err != nil {
		return GetWebinarResp{}, err
	}

	resp := GetWebinarResp{
		ID: webinar.ID,
		Author: Author{
			ID:      webinar.Author.ID,
			Name:    webinar.Author.Name,
			Picture: webinar.Author.Picture,
		},
		Picture:     webinar.Picture,
		Title:       webinar.Title,
		Description: webinar.Description,
		Category:    webinar.Category,
		StartDate:   webinar.StartDate,
		LastDate:    webinar.LastDate,
		Time:        webinar.Time,
		CreatedAt:   webinar.CreatedAt,
		UpdatedAt:   webinar.UpdatedAt,
	}

	return resp, nil
}

func (s *service) GetByCategory(ctx context.Context, category string) ([]GetWebinarResp, error) {
	webinars, err := s.webinarRepo.FindByCategory(ctx, category)
	if err != nil {
		return make([]GetWebinarResp, 0), err
	}

	resp := make([]GetWebinarResp, 0)
	for _, webinar := range webinars {
		resp = append(resp, GetWebinarResp{
			ID: webinar.ID,
			Author: Author{
				ID:      webinar.Author.ID,
				Name:    webinar.Author.Name,
				Picture: webinar.Author.Picture,
			},
			Picture:     webinar.Picture,
			Title:       webinar.Title,
			Description: webinar.Description,
			Category:    webinar.Category,
			StartDate:   webinar.StartDate,
			LastDate:    webinar.LastDate,
			Time:        webinar.Time,
			CreatedAt:   webinar.CreatedAt,
			UpdatedAt:   webinar.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Get(ctx context.Context) ([]GetWebinarResp, error) {
	webinars, err := s.webinarRepo.Find(ctx)
	if err != nil {
		return make([]GetWebinarResp, 0), err
	}

	resp := make([]GetWebinarResp, 0)
	for _, webinar := range webinars {
		resp = append(resp, GetWebinarResp{
			ID: webinar.ID,
			Author: Author{
				ID:      webinar.Author.ID,
				Name:    webinar.Author.Name,
				Picture: webinar.Author.Picture,
			},
			Picture:     webinar.Picture,
			Title:       webinar.Title,
			Description: webinar.Description,
			Category:    webinar.Category,
			StartDate:   webinar.StartDate,
			LastDate:    webinar.LastDate,
			Time:        webinar.Time,
			CreatedAt:   webinar.CreatedAt,
			UpdatedAt:   webinar.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Delete(ctx context.Context, webinarID, authorID int64) error {
	webinar, err := s.webinarRepo.FindByID(ctx, webinarID)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "Webinar not found")
	} else if err != nil {
		return err
	}

	if webinar.Author.ID != authorID {
		return app.NewError(nil, app.EForbidden, "Forbidden access, webinar not found")
	}

	return s.webinarRepo.Delete(ctx, webinarID)
}
