package discussioncomment

import (
	"context"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/discussion"
	"github.com/bagus2x/recovy/models"
)

type Service interface {
	Create(ctx context.Context, req *CreateDiscussionCommentReq) (CreateDiscussionCommentResp, error)
	GetByID(ctx context.Context, discussionCommentID int64) (GetDiscussionCommentResp, error)
	GetByDiscussionID(ctx context.Context, disussionCommentID int64) ([]GetDiscussionCommentResp, error)
	Get(ctx context.Context) ([]GetDiscussionCommentResp, error)
	Delete(ctx context.Context, discussionCommentID, commentatorID int64) error
}

type service struct {
	discussionCommentRepo Repository
	discussionRepo        discussion.Repository
}

func NewService(discussionCommentRepo Repository, discussionRepo discussion.Repository) Service {
	return &service{
		discussionCommentRepo: discussionCommentRepo,
		discussionRepo:        discussionRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreateDiscussionCommentReq) (CreateDiscussionCommentResp, error) {
	err := req.Validate()
	if err != nil {
		return CreateDiscussionCommentResp{}, err
	}

	discussion, err := s.discussionRepo.FindByID(ctx, req.DiscussionID)
	if app.ErrorCode(err) == app.ENotFound || discussion.ID == 0 {
		return CreateDiscussionCommentResp{}, app.NewError(err, app.ENotFound, "Discussion not found")
	} else if err != nil {
		return CreateDiscussionCommentResp{}, err
	}

	discussionComment := models.DiscussionComment{
		DiscussionID: req.DiscussionID,
		Commentator: models.User{
			ID: req.CommentatorID,
		},
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = s.discussionCommentRepo.Create(ctx, &discussionComment)
	if err != nil {
		return CreateDiscussionCommentResp{}, err
	}

	res := CreateDiscussionCommentResp{
		ID:            discussionComment.ID,
		DiscussionID:  discussionComment.DiscussionID,
		CommentatorID: discussionComment.Commentator.ID,
		Description:   discussionComment.Description,
		CreatedAt:     discussionComment.CreatedAt,
		UpdatedAt:     discussionComment.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, discussionCommentID int64) (GetDiscussionCommentResp, error) {
	discussionComment, err := s.discussionCommentRepo.FindByID(ctx, discussionCommentID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetDiscussionCommentResp{}, app.NewError(err, app.ENotFound, "DiscussionComment not found")
	} else if err != nil {
		return GetDiscussionCommentResp{}, err
	}

	resp := GetDiscussionCommentResp{
		ID:           discussionComment.ID,
		DiscussionID: discussionComment.DiscussionID,
		Commentator: Commentator{
			ID:      discussionComment.Commentator.ID,
			Name:    discussionComment.Commentator.Name,
			Picture: discussionComment.Commentator.Picture,
		},
		Description: discussionComment.Description,
		CreatedAt:   discussionComment.CreatedAt,
		UpdatedAt:   discussionComment.UpdatedAt,
	}

	return resp, nil
}

func (s *service) GetByDiscussionID(ctx context.Context, disussionCommentID int64) ([]GetDiscussionCommentResp, error) {
	discussionComments, err := s.discussionCommentRepo.FindByDiscussionID(ctx, disussionCommentID)
	if err != nil {
		return make([]GetDiscussionCommentResp, 0), err
	}

	resp := make([]GetDiscussionCommentResp, 0)
	for _, discussionComment := range discussionComments {
		resp = append(resp, GetDiscussionCommentResp{
			ID:           discussionComment.ID,
			DiscussionID: discussionComment.DiscussionID,
			Commentator: Commentator{
				ID:      discussionComment.Commentator.ID,
				Name:    discussionComment.Commentator.Name,
				Picture: discussionComment.Commentator.Picture,
			},
			Description: discussionComment.Description,
			CreatedAt:   discussionComment.CreatedAt,
			UpdatedAt:   discussionComment.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Get(ctx context.Context) ([]GetDiscussionCommentResp, error) {
	discussionComments, err := s.discussionCommentRepo.Find(ctx)
	if err != nil {
		return make([]GetDiscussionCommentResp, 0), err
	}

	resp := make([]GetDiscussionCommentResp, 0)
	for _, discussionComment := range discussionComments {
		resp = append(resp, GetDiscussionCommentResp{
			ID:           discussionComment.ID,
			DiscussionID: discussionComment.DiscussionID,
			Commentator: Commentator{
				ID:      discussionComment.Commentator.ID,
				Name:    discussionComment.Commentator.Name,
				Picture: discussionComment.Commentator.Picture,
			},
			Description: discussionComment.Description,
			CreatedAt:   discussionComment.CreatedAt,
			UpdatedAt:   discussionComment.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *service) Delete(ctx context.Context, discussionCommentID, commentatorID int64) error {
	discussionComment, err := s.discussionCommentRepo.FindByID(ctx, discussionCommentID)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "Discussion Comment not found")
	} else if err != nil {
		return err
	}

	if discussionComment.Commentator.ID != commentatorID {
		return app.NewError(nil, app.EForbidden, "Forbidden access, comment not found")
	}

	return s.discussionCommentRepo.Delete(ctx, discussionCommentID)
}
