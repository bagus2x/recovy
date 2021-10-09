package webinar

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
)

type CreateWebinarReq struct {
	AuthorID    int64  `json:"authorID" validate:"required,gt=0"`
	Picture     string `json:"picture" validate:"lte=512"`
	Title       string `json:"title" validate:"required,gte=5,lte=255"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category" validate:"required"`
	StartDate   int64  `json:"startDate" validate:"required,gt=0"`
	LastDate    int64  `json:"lastDate" validate:"required,gt=0"`
	Time        string `json:"time" validate:"required"`
}

func (r *CreateWebinarReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateWebinarResp struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"authorID"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	StartDate   int64  `json:"startDate"`
	LastDate    int64  `json:"lastDate"`
	Time        string `json:"time"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type Author struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GetWebinarResp struct {
	ID          int64  `json:"id"`
	Author      Author `json:"author"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	StartDate   int64  `json:"startDate"`
	LastDate    int64  `json:"lastDate"`
	Time        string `json:"time"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}
