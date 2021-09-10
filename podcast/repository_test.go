package podcast

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/bagus2x/recovy/config"
	"github.com/bagus2x/recovy/db"
	"github.com/bagus2x/recovy/models"
	"github.com/stretchr/testify/assert"
)

func getDbTest() *sql.DB {
	cfg := config.NewTest()
	db, err := db.OpenPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestRepoCreatePodcast(t *testing.T) {
	repo := NewRepository(getDbTest())
	podcast := models.Podcast{
		Author:      models.User{ID: 1},
		Picture:     "this is aurl2",
		Title:       "This is title",
		Description: "this is description",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		File:        "this is file",
	}

	err := repo.Create(context.Background(), &podcast)
	assert.NoError(t, err)
	assert.NotZero(t, podcast.ID)
}

func TestRepoFindByID(t *testing.T) {
	repo := NewRepository(getDbTest())
	p, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	t.Log(p)
}

func TestRepoFind(t *testing.T) {
	repo := NewRepository(getDbTest())
	podcasts, cursor, err := repo.Find(context.Background(), &Params{})

	assert.NoError(t, err)
	assert.NotZero(t, len(podcasts))
	t.Log(podcasts)
	t.Log(cursor)
}
