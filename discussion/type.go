package discussion

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
)

type CreateDiscussionReq struct {
	AuthorID    int64  `json:"authorID" validate:"required,gt=0"`
	Picture     string `json:"picture" validate:"lte=512"`
	Title       string `json:"title" validate:"required,lte=128"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category"  validate:"required,lte=128"`
}

func (r *CreateDiscussionReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateDiscussionResp struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"authorID"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type Author struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GetDiscussionResp struct {
	ID          int64  `json:"id"`
	Author      Author `json:"author"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}
