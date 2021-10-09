package webinar

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, webinar *models.Webinar) error
	FindByID(ctx context.Context, webinarID int64) (models.Webinar, error)
	Find(ctx context.Context) ([]models.Webinar, error)
	FindByCategory(ctx context.Context, category string) ([]models.Webinar, error)
	Delete(ctx context.Context, webinarID int64) error
}

type repository struct {
	db sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: *db,
	}
}

var create = `
	INSERT INTO
		Webinar
		(author_id, picture, title, description, category, start_date, last_date, time, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING
		id
`

func (r *repository) Create(ctx context.Context, webinar *models.Webinar) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		webinar.Author.ID,
		webinar.Picture,
		webinar.Title,
		webinar.Description,
		webinar.Category,
		webinar.StartDate,
		webinar.LastDate,
		webinar.Time,
		webinar.CreatedAt,
		webinar.UpdatedAt,
	).Scan(&webinar.ID)

	return err
}

var findByID = `
	SELECT
		w.id, au.id, au.name, au.picture, w.picture, w.title, w.description, w.category, w.start_date, w.last_date, w.time, w.created_at, w.updated_at
	FROM
		Webinar w
	JOIN
		App_User au
	ON
		w.author_id = au.id
	WHERE
		w.id = $1
`

func (r *repository) FindByID(ctx context.Context, webinarID int64) (models.Webinar, error) {
	var webinar models.Webinar

	err := r.db.QueryRowContext(ctx, findByID, webinarID).Scan(
		&webinar.ID,
		&webinar.Author.ID,
		&webinar.Author.Name,
		&webinar.Author.Picture,
		&webinar.Picture,
		&webinar.Title,
		&webinar.Description,
		&webinar.Category,
		&webinar.StartDate,
		&webinar.LastDate,
		&webinar.Time,
		&webinar.CreatedAt,
		&webinar.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Webinar{}, app.NewError(err, app.ENotFound)
	}

	return webinar, nil
}

var find = `
	SELECT
		w.id, au.id, au.name, au.picture, w.picture, w.title, w.description, w.category, w.start_date, w.last_date, w.time, w.created_at, w.updated_at
	FROM
		Webinar w
	JOIN
		App_User au
	ON
		w.author_id = au.id
	LIMIT
		100
`

func (r *repository) Find(ctx context.Context) ([]models.Webinar, error) {
	webinars := make([]models.Webinar, 0)

	rows, err := r.db.QueryContext(ctx, find)
	if err != nil {
		return webinars, err
	}

	for rows.Next() {
		var webinar models.Webinar

		err := rows.Scan(
			&webinar.ID,
			&webinar.Author.ID,
			&webinar.Author.Name,
			&webinar.Author.Picture,
			&webinar.Picture,
			&webinar.Title,
			&webinar.Description,
			&webinar.Category,
			&webinar.StartDate,
			&webinar.LastDate,
			&webinar.Time,
			&webinar.CreatedAt,
			&webinar.UpdatedAt,
		)
		if err != nil {
			return webinars, err
		}

		webinars = append(webinars, webinar)
	}
	if rows.Err() != nil {
		return webinars, err
	}

	return webinars, nil
}

var findByCategory = `
	SELECT
		w.id, au.id, au.name, au.picture, w.picture, w.title, w.description, w.category, w.start_date, w.last_date, w.time, w.created_at, w.updated_at
	FROM
		Webinar w
	JOIN
		App_User au
	ON
		w.author_id = au.id
	WHERE
		w.category = $1
	LIMIT
		100
`

func (r *repository) FindByCategory(ctx context.Context, category string) ([]models.Webinar, error) {
	webinars := make([]models.Webinar, 0)

	rows, err := r.db.QueryContext(ctx, findByCategory, category)
	if err != nil {
		return webinars, err
	}

	for rows.Next() {
		var webinar models.Webinar

		err := rows.Scan(
			&webinar.ID,
			&webinar.Author.ID,
			&webinar.Author.Name,
			&webinar.Author.Picture,
			&webinar.Picture,
			&webinar.Title,
			&webinar.Description,
			&webinar.Category,
			&webinar.StartDate,
			&webinar.LastDate,
			&webinar.Time,
			&webinar.CreatedAt,
			&webinar.UpdatedAt,
		)
		if err != nil {
			return webinars, err
		}

		webinars = append(webinars, webinar)
	}
	if rows.Err() != nil {
		return webinars, err
	}

	return webinars, nil
}

var delete = `
	DELETE FROM
		Webinar
	WHERE
		id = $1
`

func (r *repository) Delete(ctx context.Context, webinarID int64) error {
	_, err := r.db.ExecContext(ctx, delete, webinarID)

	return err
}
