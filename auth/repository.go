package auth

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, userID int64) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

var create = `
	INSERT INTO 
		App_User
		(name, email, picture, password, created_at, updated_at, deleted_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id
`

func (r *repository) Create(ctx context.Context, user *models.User) error {
	err := r.db.QueryRowxContext(
		ctx,
		create,
		user.Name,
		user.Email,
		user.Picture,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	).Scan(&user.ID)

	return err
}

var findByID = `
	SELECT
		id, name, email, picture, password, created_at, updated_at, deleted_at
	FROM
		App_User
	WHERE
		id = $1
`

func (r *repository) FindByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, findByID, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

var findByEmail = `
	SELECT
		id, name, email, picture, password, created_at, updated_at, deleted_at
	FROM
		App_User
	WHERE
		email = $1
`

func (r *repository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, findByEmail, email)
	if err != nil && err == sql.ErrNoRows {
		return models.User{}, app.NewError(err, app.ENotFound)
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
