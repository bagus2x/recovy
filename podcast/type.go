package podcast

import (
	"github.com/bagus2x/recovy/app"
	"github.com/go-playground/validator/v10"
)

type Podcast struct {
	ID          int64  `json:"id"`
	Author      Author `json:"author"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file"`
	Starred     bool   `json:"starred"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type Author struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Params struct {
	Cursor    int64
	Limit     int64
	Direction string
}

type Cursor struct {
	Next     int64 `json:"next"`
	Previous int64 `json:"previous"`
}

type CreatePodcastReq struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"authorID" validate:"required"`
	Picture     string `json:"picture" validate:"lte=512"`
	Title       string `json:"title" validate:"required,gte=5,lte=255"`
	Description string `json:"description" validate:"lte=512"`
	File        string `json:"file" validate:"required,gte=5,lte=255"`
}

func (r *CreatePodcastReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreatePodcastResp struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"authorID"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type GetPodcastByIDResp Podcast

type GetPodcastsResp struct {
	Cursor   Cursor    `json:"cursor"`
	Podcasts []Podcast `json:"podcasts"`
}

type StarPodcastReq struct {
	PodcastID int64 `json:"podcast_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
}

func (r *StarPodcastReq) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type UnstarPodcastReq StarPodcastReq
