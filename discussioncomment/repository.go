package discussioncomment

import (
	"context"
	"database/sql"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/models"
)

type Repository interface {
	Create(ctx context.Context, discussionComment *models.DiscussionComment) error
	FindByID(ctx context.Context, discussionCommentID int64) (models.DiscussionComment, error)
	Find(ctx context.Context) ([]models.DiscussionComment, error)
	FindByDiscussionID(ctx context.Context, discussionID int64) ([]models.DiscussionComment, error)
	Delete(ctx context.Context, discussionCommentID int64) error
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
		Discussion_Comment
		(discussion_id, commentator_id,  description, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		id
`

func (r *repository) Create(ctx context.Context, discussionComment *models.DiscussionComment) error {
	err := r.db.QueryRowContext(
		ctx,
		create,
		discussionComment.DiscussionID,
		discussionComment.Commentator.ID,
		discussionComment.Description,
		discussionComment.CreatedAt,
		discussionComment.UpdatedAt,
	).Scan(&discussionComment.ID)

	return err
}

var findByID = `
	SELECT
		dc.id, dc.discussion_id, au.id, au.name, au.picture,  dc.description, dc.created_at, dc.updated_at
	FROM
		Discussion_Comment dc
	JOIN
		App_User au
	ON
		dc.commentator_id = au.id
	WHERE
		dc.id = $1
`

func (r *repository) FindByID(ctx context.Context, discussionCommentID int64) (models.DiscussionComment, error) {
	var discussionComment models.DiscussionComment

	err := r.db.QueryRowContext(ctx, findByID, discussionCommentID).Scan(
		&discussionComment.ID,
		&discussionComment.DiscussionID,
		&discussionComment.Commentator.ID,
		&discussionComment.Commentator.Name,
		&discussionComment.Commentator.Picture,
		&discussionComment.Description,
		&discussionComment.CreatedAt,
		&discussionComment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.DiscussionComment{}, app.NewError(err, app.ENotFound)
	}

	return discussionComment, nil
}

var find = `
	SELECT
		dc.id, dc.discussion_id, au.id, au.name, au.picture,  dc.description, dc.created_at, dc.updated_at
	FROM
		Discussion_Comment dc
	JOIN
		App_User au
	ON
		dc.commentator_id = au.id
	LIMIT
		100
`

func (r *repository) Find(ctx context.Context) ([]models.DiscussionComment, error) {
	discussionComments := make([]models.DiscussionComment, 0)

	rows, err := r.db.QueryContext(ctx, find)
	if err != nil {
		return discussionComments, err
	}

	for rows.Next() {
		var discussionComment models.DiscussionComment

		err := rows.Scan(
			&discussionComment.ID,
			&discussionComment.DiscussionID,
			&discussionComment.Commentator.ID,
			&discussionComment.Commentator.Name,
			&discussionComment.Commentator.Picture,
			&discussionComment.Description,
			&discussionComment.CreatedAt,
			&discussionComment.UpdatedAt,
		)
		if err != nil {
			return discussionComments, err
		}

		discussionComments = append(discussionComments, discussionComment)
	}
	if rows.Err() != nil {
		return discussionComments, err
	}

	return discussionComments, nil
}

var findByDiscussionID = `
	SELECT
		dc.id, dc.discussion_id, au.id, au.name, au.picture,  dc.description, dc.created_at, dc.updated_at
	FROM
		Discussion_Comment dc
	JOIN
		App_User au
	ON
		dc.commentator_id = au.id
	WHERE
		dc.discussion_id = $1
	LIMIT
		100
`

func (r *repository) FindByDiscussionID(ctx context.Context, discussionCommentID int64) ([]models.DiscussionComment, error) {
	discussionComments := make([]models.DiscussionComment, 0)

	rows, err := r.db.QueryContext(ctx, findByDiscussionID, discussionCommentID)
	if err != nil {
		return discussionComments, err
	}

	for rows.Next() {
		var discussionComment models.DiscussionComment

		err := rows.Scan(
			&discussionComment.ID,
			&discussionComment.DiscussionID,
			&discussionComment.Commentator.ID,
			&discussionComment.Commentator.Name,
			&discussionComment.Commentator.Picture,
			&discussionComment.Description,
			&discussionComment.CreatedAt,
			&discussionComment.UpdatedAt,
		)
		if err != nil {
			return discussionComments, err
		}

		discussionComments = append(discussionComments, discussionComment)
	}
	if rows.Err() != nil {
		return discussionComments, err
	}

	return discussionComments, nil
}

var delete = `
	DELETE FROM
		Discussion_Comment
	WHERE
		id = $1
`

func (r *repository) Delete(ctx context.Context, discussionCommentID int64) error {
	_, err := r.db.ExecContext(ctx, delete, discussionCommentID)

	return err
}
