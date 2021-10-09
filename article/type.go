package article

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
)

type CreateArticleReq struct {
	AuthorID    int64  `json:"author_id" validate:"required,gt=0"`
	Picture     string `json:"picture" validate:"lte=512"`
	Title       string `json:"title" validate:"required,lte=255"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category" validate:"required,lte=128"`
}

func (r *CreateArticleReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateArticleResp struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"author_id"`
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

type GetArticleResp struct {
	ID          int64  `json:"id"`
	Author      Author `json:"author"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}
