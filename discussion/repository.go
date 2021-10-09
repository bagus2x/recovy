package discussion

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, discussion *models.Discussion) error
	FindByID(ctx context.Context, discussionID int64) (models.Discussion, error)
	Find(ctx context.Context) ([]models.Discussion, error)
	FindByCategory(ctx context.Context, category string) ([]models.Discussion, error)
	Delete(ctx context.Context, discussionID int64) error
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
		Discussion
		(author_id, picture, title, description, category, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id
`

func (r *repository) Create(ctx context.Context, discussion *models.Discussion) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		discussion.Author.ID,
		discussion.Picture,
		discussion.Title,
		discussion.Description,
		discussion.Category,
		discussion.CreatedAt,
		discussion.UpdatedAt,
	).Scan(&discussion.ID)

	return err
}

var findByID = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Discussion a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	WHERE
		a.id = $1
`

func (r *repository) FindByID(ctx context.Context, discussionID int64) (models.Discussion, error) {
	var discussion models.Discussion

	err := r.db.QueryRowContext(ctx, findByID, discussionID).Scan(
		&discussion.ID,
		&discussion.Author.ID,
		&discussion.Author.Name,
		&discussion.Author.Picture,
		&discussion.Picture,
		&discussion.Title,
		&discussion.Description,
		&discussion.Category,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Discussion{}, app.NewError(err, app.ENotFound)
	}

	return discussion, nil
}

var find = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Discussion a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	LIMIT
		100
`

func (r *repository) Find(ctx context.Context) ([]models.Discussion, error) {
	discussions := make([]models.Discussion, 0)

	rows, err := r.db.QueryContext(ctx, find)
	if err != nil {
		return discussions, err
	}

	for rows.Next() {
		var discussion models.Discussion

		err := rows.Scan(
			&discussion.ID,
			&discussion.Author.ID,
			&discussion.Author.Name,
			&discussion.Author.Picture,
			&discussion.Picture,
			&discussion.Title,
			&discussion.Description,
			&discussion.Category,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
		)
		if err != nil {
			return discussions, err
		}

		discussions = append(discussions, discussion)
	}
	if rows.Err() != nil {
		return discussions, err
	}

	return discussions, nil
}

var findByCategory = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Discussion a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	WHERE
		a.category = $1
	LIMIT
		100
`

func (r *repository) FindByCategory(ctx context.Context, category string) ([]models.Discussion, error) {
	discussions := make([]models.Discussion, 0)

	rows, err := r.db.QueryContext(ctx, findByCategory, category)
	if err != nil {
		return discussions, err
	}

	for rows.Next() {
		var discussion models.Discussion

		err := rows.Scan(
			&discussion.ID,
			&discussion.Author.ID,
			&discussion.Author.Name,
			&discussion.Author.Picture,
			&discussion.Picture,
			&discussion.Title,
			&discussion.Description,
			&discussion.Category,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
		)
		if err != nil {
			return discussions, err
		}

		discussions = append(discussions, discussion)
	}
	if rows.Err() != nil {
		return discussions, err
	}

	return discussions, nil
}

var delete = `
	DELETE FROM
		Discussion
	WHERE
		id = $1
`

func (r *repository) Delete(ctx context.Context, discussionID int64) error {
	_, err := r.db.ExecContext(ctx, delete, discussionID)

	return err
}
