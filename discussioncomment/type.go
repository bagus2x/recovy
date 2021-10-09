package discussioncomment

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
)

type CreateDiscussionCommentReq struct {
	DiscussionID  int64  `json:"discussionID" validate:"required,gt=0"`
	CommentatorID int64  `json:"commentatorID" validate:"required,gt=0" `
	Description   string `json:"description" validate:"required"`
}

func (r *CreateDiscussionCommentReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateDiscussionCommentResp struct {
	ID            int64  `json:"id"`
	DiscussionID  int64  `json:"discussion_id"`
	CommentatorID int64  `json:"commentator_id"`
	Description   string `json:"description"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type Commentator struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GetDiscussionCommentResp struct {
	ID           int64       `json:"id"`
	DiscussionID int64       `json:"discussionID"`
	Commentator  Commentator `json:"commentator"`
	Description  string      `json:"description"`
	CreatedAt    int64       `json:"createdAt"`
	UpdatedAt    int64       `json:"updatedAt"`
}
