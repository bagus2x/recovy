package discussion

import (
	"context"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Service interface {
	Create(ctx context.Context, req *CreateDiscussionReq) (CreateDiscussionResp, error)
	GetByID(ctx context.Context, discussionID int64) (GetDiscussionResp, error)
	GetByCategory(ctx context.Context, category string) ([]GetDiscussionResp, error)
	Get(ctx context.Context) ([]GetDiscussionResp, error)
	Delete(ctx context.Context, discussionID, authorID int64) error
}

type service struct {
	discussionRepo Repository
}

func NewService(discussionRepo Repository) Service {
	return &service{
		discussionRepo: discussionRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreateDiscussionReq) (CreateDiscussionResp, error) {
	err := req.Validate()
	if err != nil {
		return CreateDiscussionResp{}, err
	}

	discussion := models.Discussion{
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

	err = s.discussionRepo.Create(ctx, &discussion)
	if err != nil {
		return CreateDiscussionResp{}, err
	}

	res := CreateDiscussionResp{
		ID:          discussion.ID,
		AuthorID:    discussion.Author.ID,
		Picture:     discussion.Picture,
		Title:       discussion.Title,
		Description: discussion.Description,
		Category:    discussion.Category,
		CreatedAt:   discussion.CreatedAt,
		UpdatedAt:   discussion.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, discussionID int64) (GetDiscussionResp, error) {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetDiscussionResp{}, app.NewError(err, app.ENotFound, "Discussion not found")
	} else if err != nil {
		return GetDiscussionResp{}, err
	}

	resp := GetDiscussionResp{
		ID: discussion.ID,
		Author: Author{
			ID:      discussion.Author.ID,
			Name:    discussion.Author.Name,
			Picture: discussion.Author.Picture,
		},
		Picture:     discussion.Picture,
		Title:       discussion.Title,
		Description: discussion.Description,
		Category:    discussion.Category,
		CreatedAt:   discussion.CreatedAt,
		UpdatedAt:   discussion.UpdatedAt,
	}

	return resp, nil
}

func (s *service) GetByCategory(ctx context.Context, category string) ([]GetDiscussionResp, error) {
	discussions, err := s.discussionRepo.FindByCategory(ctx, category)
	if err != nil {
		return make([]GetDiscussionResp, 0), err
	}

	resp := make([]GetDiscussionResp, 0)
	for _, discussion := range discussions {
		resp = append(resp, GetDiscussionResp{
			ID: discussion.ID,
			Author: Author{
				ID:      discussion.Author.ID,
				Name:    discussion.Author.Name,
				Picture: discussion.Author.Picture,
			},
			Picture:     discussion.Picture,
			Title:       discussion.Title,
			Description: discussion.Description,
			Category:    discussion.Category,
			CreatedAt:   discussion.CreatedAt,
			UpdatedAt:   discussion.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Get(ctx context.Context) ([]GetDiscussionResp, error) {
	discussions, err := s.discussionRepo.Find(ctx)
	if err != nil {
		return make([]GetDiscussionResp, 0), err
	}

	resp := make([]GetDiscussionResp, 0)
	for _, discussion := range discussions {
		resp = append(resp, GetDiscussionResp{
			ID: discussion.ID,
			Author: Author{
				ID:      discussion.Author.ID,
				Name:    discussion.Author.Name,
				Picture: discussion.Author.Picture,
			},
			Picture:     discussion.Picture,
			Title:       discussion.Title,
			Description: discussion.Description,
			Category:    discussion.Category,
			CreatedAt:   discussion.CreatedAt,
			UpdatedAt:   discussion.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Delete(ctx context.Context, discussionID, authorID int64) error {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "Discussion not found")
	} else if err != nil {
		return err
	}

	if discussion.Author.ID != authorID {
		return app.NewError(nil, app.EForbidden, "Forbidden access, discussion not found")
	}

	return s.discussionRepo.Delete(ctx, discussionID)
}
