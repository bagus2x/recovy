package podcast

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, podcast *models.Podcast) error
	FindByID(ctx context.Context, podcastID int64) (models.Podcast, error)
	Find(ctx context.Context, params *Params) ([]models.Podcast, Cursor, error)
	Delete(ctx context.Context, podcastID int64) error
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
				Podcast
				(author_id, picture, title, description, file, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
			RETURNING
				id
`

func (r *repository) Create(ctx context.Context, podcast *models.Podcast) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		podcast.Author.ID,
		podcast.Picture,
		podcast.Title,
		podcast.Description,
		podcast.File,
		podcast.CreatedAt,
		podcast.UpdatedAt,
	).Scan(&podcast.ID)

	return err
}

var findByID = `
			SELECT
				pt.id, au.id, au.name, au.picture, pt.picture, pt.title, pt.description, pt.file, pt.created_at, pt.updated_at
			FROM
				Podcast pt
			JOIN
				App_User au
			ON
				pt.author_id = au.id
			WHERE
				pt.id = $1
`

func (r *repository) FindByID(ctx context.Context, podcastID int64) (models.Podcast, error) {
	var podcast models.Podcast

	err := r.db.QueryRowContext(ctx, findByID, podcastID).Scan(
		&podcast.ID,
		&podcast.Author.ID,
		&podcast.Author.Name,
		&podcast.Author.Picture,
		&podcast.Picture,
		&podcast.Title,
		&podcast.Description,
		&podcast.File,
		&podcast.CreatedAt,
		&podcast.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Podcast{}, app.NewError(err, app.ENotFound)
	} else if err != nil {
		return models.Podcast{}, err
	}

	return podcast, nil
}

func findNext(params *Params) string {
	podcasts := strings.Builder{}
	podcasts.WriteString(`
		SELECT
			pt.id, au.id, au.name, au.picture, pt.picture, pt.title, pt.description, pt.file, pt.created_at, pt.updated_at
		FROM
			Podcast pt
		JOIN
			App_User au
		ON
			pt.author_id = au.id
		WHERE`,
	)

	// Cursor
	cursor := params.Cursor
	if cursor == 0 {
		cursor = math.MaxInt32
	}

	limit := params.Limit
	if limit == 0 {
		limit = 10
	}

	fmt.Fprintf(&podcasts, " pt.id < %d ORDER BY pt.id DESC LIMIT %d ", cursor, limit)

	return podcasts.String()
}

func findPrev(params *Params) string {
	podcasts := strings.Builder{}
	podcasts.WriteString(`
		WITH prev_mode AS (
			SELECT
				pt.id as podcast_id, au.id, au.name, au.picture, pt.picture, pt.title, pt.description, pt.file, pt.created_at, pt.updated_at
			FROM
				Podcast pt
			JOIN
				App_User au
			ON
				pt.author_id = au.id
			WHERE`,
	)

	// Cursor
	limit := params.Limit
	if limit == 0 {
		limit = 10
	}

	fmt.Fprintf(&podcasts, " pt.id > %d LIMIT %d ", params.Cursor, limit)
	podcasts.WriteString(") SELECT * FROM prev_mode ORDER BY podcast_id DESC")

	return podcasts.String()
}

func find(params *Params) string {
	if params.Direction == "previous" {
		return findPrev(params)
	}

	return findNext(params)
}

func (r *repository) Find(ctx context.Context, params *Params) ([]models.Podcast, Cursor, error) {
	rows, err := r.db.QueryContext(ctx, find(params))
	if err != nil {
		return nil, Cursor{}, err
	}
	defer rows.Close()

	podcasts := make([]models.Podcast, 0)

	for rows.Next() {
		var podcast models.Podcast

		err := rows.Scan(
			&podcast.ID,
			&podcast.Author.ID,
			&podcast.Author.Name,
			&podcast.Author.Picture,
			&podcast.Picture,
			&podcast.Title,
			&podcast.Description,
			&podcast.File,
			&podcast.CreatedAt,
			&podcast.UpdatedAt,
		)
		if err != nil {
			return nil, Cursor{}, err
		}

		podcasts = append(podcasts, podcast)
	}
	if rows.Err() != nil {
		return nil, Cursor{}, err
	}

	var cursor Cursor
	if len(podcasts) > 0 {
		cursor.Next = podcasts[len(podcasts)-1].ID
		cursor.Previous = podcasts[0].ID
	}

	return podcasts, cursor, nil
}

var delete = `
			DELETE FROM
				Podcast
			WHERE
				id = $1
`

func (r *repository) Delete(ctx context.Context, podcastID int64) error {
	res, err := r.db.ExecContext(ctx, delete, podcastID)
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
