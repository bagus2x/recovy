package podcast

import (
	"context"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
	"github.com/bagus2x/recovy/starredpodcast"
)

type Service interface {
	Create(ctx context.Context, req *CreatePodcastReq) (CreatePodcastResp, error)
	GetByID(ctx context.Context, podcastID int64) (GetPodcastByIDResp, error)
	GetByParams(ctx context.Context, params *Params) (GetPodcastsResp, error)
	StarPodcast(ctx context.Context, req *StarPodcastReq) error
	UnstarPodcast(ctx context.Context, req *StarPodcastReq) error
	DeleteByPodcastIDAndAthorID(ctx context.Context, podcastID, authorID int64) error
}

type service struct {
	podcastRepo        Repository
	starredPodcastRepo starredpodcast.Repository
}

func NewService(podcastRepo Repository, starredPodcastRepo starredpodcast.Repository) Service {
	return &service{
		podcastRepo:        podcastRepo,
		starredPodcastRepo: starredPodcastRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreatePodcastReq) (CreatePodcastResp, error) {
	err := req.Validate()
	if err != nil {
		return CreatePodcastResp{}, err
	}

	podcast := models.Podcast{
		Author:      models.User{ID: req.AuthorID},
		Picture:     req.Picture,
		Title:       req.Title,
		Description: req.Description,
		File:        req.File,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	err = s.podcastRepo.Create(ctx, &podcast)
	if err != nil {
		return CreatePodcastResp{}, err
	}

	res := CreatePodcastResp{
		ID:          podcast.ID,
		AuthorID:    podcast.Author.ID,
		Picture:     podcast.Picture,
		Title:       podcast.Title,
		Description: podcast.Description,
		File:        podcast.File,
		CreatedAt:   podcast.CreatedAt,
		UpdatedAt:   podcast.UpdatedAt,
	}

	return res, nil
}

func (s *service) StarPodcast(ctx context.Context, req *StarPodcastReq) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	return s.starredPodcastRepo.Create(ctx, &models.StarredPodcast{
		Podcast: models.Podcast{
			ID: req.PodcastID,
		},
		User: models.User{
			ID: req.UserID,
		},
	})
}

func (s *service) UnstarPodcast(ctx context.Context, req *StarPodcastReq) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	return s.starredPodcastRepo.DeleteByPodcastIDAndUserID(ctx, &models.StarredPodcast{
		Podcast: models.Podcast{
			ID: req.PodcastID,
		},
		User: models.User{
			ID: req.UserID,
		},
	})
}

func (s *service) GetByID(ctx context.Context, podcastID int64) (GetPodcastByIDResp, error) {
	podcast, err := s.podcastRepo.FindByID(ctx, podcastID)
	if app.ErrorCode(err) == app.ENotFound {
		return GetPodcastByIDResp{}, app.NewError(err, app.ENotFound, "Podcast not found")
	} else if err != nil {
		return GetPodcastByIDResp{}, nil
	}

	res := GetPodcastByIDResp{
		ID: podcastID,
		Author: Author{
			ID:      podcast.Author.ID,
			Name:    podcast.Author.Name,
			Picture: podcast.Author.Picture,
		},
		Picture:     podcast.Picture,
		Title:       podcast.Title,
		Description: podcast.Description,
		File:        podcast.File,
		CreatedAt:   podcast.CreatedAt,
		UpdatedAt:   podcast.UpdatedAt,
	}

	userID, ok := ctx.Value("userID").(int64)
	if ok {
		star, err := s.starredPodcastRepo.FindByPodcastIDAndUserID(ctx, podcastID, userID)
		if err != nil && app.ErrorCode(err) != app.ENotFound {
			return GetPodcastByIDResp{}, nil
		}

		res.Starred = star != models.StarredPodcast{}
	}

	return res, nil
}

func (s *service) GetByParams(ctx context.Context, params *Params) (GetPodcastsResp, error) {
	podcasts, cursor, err := s.podcastRepo.Find(ctx, params)
	if err != nil {
		return GetPodcastsResp{}, err
	}

	res := GetPodcastsResp{
		Cursor: cursor,
	}

	userID, ok := ctx.Value("userID").(int64)

	for _, podcast := range podcasts {
		var starred bool
		if ok {
			star, err := s.starredPodcastRepo.FindByPodcastIDAndUserID(ctx, podcast.ID, userID)
			if err != nil && app.ErrorCode(err) != app.ENotFound {
				return GetPodcastsResp{}, err
			}
			starred = star != models.StarredPodcast{}
		}

		res.Podcasts = append(res.Podcasts, Podcast{
			ID: podcast.ID,
			Author: Author{
				ID:      podcast.Author.ID,
				Name:    podcast.Author.Name,
				Picture: podcast.Author.Picture,
			},
			Picture:     podcast.Picture,
			Title:       podcast.Description,
			Description: podcast.Description,
			File:        podcast.File,
			Starred:     starred,
			CreatedAt:   podcast.CreatedAt,
			UpdatedAt:   podcast.UpdatedAt,
		})

	}

	return res, nil
}

func (s *service) DeleteByPodcastIDAndAthorID(ctx context.Context, podcastID, authorID int64) error {
	podcast, err := s.podcastRepo.FindByID(ctx, podcastID)
	if app.ErrorCode(err) == app.ENotFound || podcast == (models.Podcast{}) {
		return app.NewError(err, app.ENotFound, "Podcast not found")
	} else if err != nil {
		return nil
	}

	if podcast.Author.ID != authorID {
		return app.NewError(nil, app.EForbidden)
	}

	err = s.podcastRepo.Delete(ctx, podcastID)
	if err != nil {
		return nil
	}

	return nil
}
