package article

import (
	"context"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Service interface {
	Create(ctx context.Context, req *CreateArticleReq) (CreateArticleResp, error)
	GetByID(ctx context.Context, articleID int64) (GetArticleResp, error)
	GetByCategory(ctx context.Context, category string) ([]GetArticleResp, error)
	Get(ctx context.Context) ([]GetArticleResp, error)
	Delete(ctx context.Context, articleID, authorID int64) error
}

type service struct {
	articleRepo Repository
}

func NewService(articleRepo Repository) Service {
	return &service{
		articleRepo: articleRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreateArticleReq) (CreateArticleResp, error) {
	err := req.Validate()
	if err != nil {
		return CreateArticleResp{}, err
	}

	article := models.Article{
		Author: models.User{
			ID: req.AuthorID,
		},
		Picture:     req.Picture,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = s.articleRepo.Create(ctx, &article)
	if err != nil {
		return CreateArticleResp{}, err
	}

	res := CreateArticleResp{
		ID:          article.ID,
		AuthorID:    article.Author.ID,
		Picture:     article.Picture,
		Title:       article.Title,
		Description: article.Description,
		Category:    article.Category,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, articleID int64) (GetArticleResp, error) {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetArticleResp{}, app.NewError(err, app.ENotFound, "Article not found")
	} else if err != nil {
		return GetArticleResp{}, err
	}

	resp := GetArticleResp{
		ID: article.ID,
		Author: Author{
			ID:      article.Author.ID,
			Name:    article.Author.Name,
			Picture: article.Author.Picture,
		},
		Picture:     article.Picture,
		Title:       article.Title,
		Description: article.Description,
		Category:    article.Category,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return resp, nil
}

func (s *service) GetByCategory(ctx context.Context, category string) ([]GetArticleResp, error) {
	articles, err := s.articleRepo.FindByCategory(ctx, category)
	if err != nil {
		return make([]GetArticleResp, 0), err
	}

	resp := make([]GetArticleResp, 0)
	for _, article := range articles {
		resp = append(resp, GetArticleResp{
			ID: article.ID,
			Author: Author{
				ID:      article.Author.ID,
				Name:    article.Author.Name,
				Picture: article.Author.Picture,
			},
			Picture:     article.Picture,
			Title:       article.Title,
			Description: article.Description,
			Category:    article.Category,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Get(ctx context.Context) ([]GetArticleResp, error) {
	articles, err := s.articleRepo.Find(ctx)
	if err != nil {
		return make([]GetArticleResp, 0), err
	}

	resp := make([]GetArticleResp, 0)
	for _, article := range articles {
		resp = append(resp, GetArticleResp{
			ID: article.ID,
			Author: Author{
				ID:      article.Author.ID,
				Name:    article.Author.Name,
				Picture: article.Author.Picture,
			},
			Picture:     article.Picture,
			Title:       article.Title,
			Description: article.Description,
			Category:    article.Category,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Delete(ctx context.Context, articleID, authorID int64) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "Discussion not found")
	} else if err != nil {
		return err
	}

	if article.Author.ID != authorID {
		return app.NewError(nil, app.EForbidden, "Forbidden access, article not found")
	}

	return s.articleRepo.Delete(ctx, articleID)
}
