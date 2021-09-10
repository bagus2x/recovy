package starredpodcast

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, starredpodcast *models.StarredPodcast) error
	FindByPodcastIDAndUserID(ctx context.Context, podcastID, userID int64) (models.StarredPodcast, error)
	DeleteByPodcastIDAndUserID(ctx context.Context, starredpodcast *models.StarredPodcast) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

var create = `
			INSERT INTO
				Starred_Podcast
				(podcast_id, user_id, created_at)
			VALUES
				($1, $2, $3)
			RETURNING
				id
`

func (r *repository) Create(ctx context.Context, starredpodcast *models.StarredPodcast) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		starredpodcast.Podcast.ID,
		starredpodcast.User.ID,
		starredpodcast.CreatedAt,
	).Scan(&starredpodcast.ID)

	return err
}

var findByID = `
			SELECT
				id, podcast_id, user_id, created_at
			FROM
				Starred_Podcast
			WHERE
				podcast_id = $1 AND user_id = $2
`

func (r *repository) FindByPodcastIDAndUserID(ctx context.Context, podcastID, userID int64) (models.StarredPodcast, error) {
	var starredpodcast models.StarredPodcast

	err := r.db.QueryRowContext(ctx, findByID, podcastID, userID).Scan(
		&starredpodcast.ID,
		&starredpodcast.Podcast.ID,
		&starredpodcast.User.ID,
		&starredpodcast.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.StarredPodcast{}, app.NewError(err, app.ENotFound)
	} else if err != nil {
		return models.StarredPodcast{}, err
	}

	return starredpodcast, nil
}

var delete = `
			DELETE FROM
				Starred_Podcast
			WHERE
				podcast_id = $1 AND user_id = $2
`

func (r *repository) DeleteByPodcastIDAndUserID(ctx context.Context, starredpodcast *models.StarredPodcast) error {
	res, err := r.db.ExecContext(ctx, delete, starredpodcast.Podcast.ID, starredpodcast.User.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return app.NewError(nil, app.ENotFound)
	}

	return nil
}
