package article

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, article *models.Article) error
	FindByID(ctx context.Context, articleID int64) (models.Article, error)
	Find(ctx context.Context) ([]models.Article, error)
	FindByCategory(ctx context.Context, category string) ([]models.Article, error)
	Delete(ctx context.Context, articleID int64) error
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
		Article
		(author_id, picture, title, description, category, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id
`

func (r *repository) Create(ctx context.Context, article *models.Article) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		article.Author.ID,
		article.Picture,
		article.Title,
		article.Description,
		article.Category,
		article.CreatedAt,
		article.UpdatedAt,
	).Scan(&article.ID)

	return err
}

var findByID = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Article a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	WHERE
		a.id = $1
`

func (r *repository) FindByID(ctx context.Context, articleID int64) (models.Article, error) {
	var article models.Article

	err := r.db.QueryRowContext(ctx, findByID, articleID).Scan(
		&article.ID,
		&article.Author.ID,
		&article.Author.Name,
		&article.Author.Picture,
		&article.Picture,
		&article.Title,
		&article.Description,
		&article.Category,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Article{}, app.NewError(err, app.ENotFound)
	}

	return article, nil
}

var find = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Article a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	LIMIT
		100
`

func (r *repository) Find(ctx context.Context) ([]models.Article, error) {
	articles := make([]models.Article, 0)

	rows, err := r.db.QueryContext(ctx, find)
	if err != nil {
		return articles, err
	}

	for rows.Next() {
		var article models.Article

		err := rows.Scan(
			&article.ID,
			&article.Author.ID,
			&article.Author.Name,
			&article.Author.Picture,
			&article.Picture,
			&article.Title,
			&article.Description,
			&article.Category,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return articles, err
		}

		articles = append(articles, article)
	}
	if rows.Err() != nil {
		return articles, err
	}

	return articles, nil
}

var findByCategory = `
	SELECT
		a.id, au.id, au.name, au.picture, a.picture, a.title, a.description, a.category, a.created_at, a.updated_at
	FROM
		Article a
	JOIN
		App_User au
	ON
		a.author_id = au.id
	WHERE
		a.category = $1
	LIMIT
		100
`

func (r *repository) FindByCategory(ctx context.Context, category string) ([]models.Article, error) {
	articles := make([]models.Article, 0)

	rows, err := r.db.QueryContext(ctx, findByCategory, category)
	if err != nil {
		return articles, err
	}

	for rows.Next() {
		var article models.Article

		err := rows.Scan(
			&article.ID,
			&article.Author.ID,
			&article.Author.Name,
			&article.Author.Picture,
			&article.Picture,
			&article.Title,
			&article.Description,
			&article.Category,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return articles, err
		}

		articles = append(articles, article)
	}
	if rows.Err() != nil {
		return articles, err
	}

	return articles, nil
}

var delete = `
	DELETE FROM
		Article
	WHERE
		id = $1
`

func (r *repository) Delete(ctx context.Context, articleID int64) error {
	_, err := r.db.ExecContext(ctx, delete, articleID)

	return err
}
