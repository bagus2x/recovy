package auth

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

func TestCreateUser(t *testing.T) {
	repo := NewRepository(getDbTest())
	user := models.User{
		Name:      "budi",
		Email:     "budi@gmail.com",
		Password:  "budi123",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err := repo.Create(context.Background(), &user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestFindByID(t *testing.T) {
	repo := NewRepository(getDbTest())

	user, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	t.Log(user)
	assert.NotZero(t, user.ID)
}

func TestFindByEmail(t *testing.T) {
	repo := NewRepository(getDbTest())

	user, err := repo.FindByEmail(context.Background(), "budi@gmail.com")
	assert.NoError(t, err)
	t.Log(user)
	assert.NotZero(t, user.ID)
}
