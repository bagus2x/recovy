package auth

import (
	"fmt"
	"time"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/config"
	"github.com/coocood/freecache"
)

type CacheRepository interface {
	SetRefreshToken(userID int64, token string) error
	GetRefreshToken(userID int64) (string, error)
	DeleteRefreshToken(userID int64) bool
}

type cacheRepository struct {
	cache *freecache.Cache
	cfg   *config.Config
}

func NewCacheRepository(cache *freecache.Cache, cfg *config.Config) CacheRepository {
	return &cacheRepository{
		cache: cache,
		cfg:   cfg,
	}
}

func (c cacheRepository) SetRefreshToken(userID int64, token string) error {
	exp := int(time.Now().Add(time.Second * time.Duration(c.cfg.RefreshTokenLifetime())).Unix())
	return c.cache.Set([]byte(fmt.Sprintf("%d", userID)), []byte(token), exp)
}

func (c cacheRepository) GetRefreshToken(userID int64) (string, error) {
	b, err := c.cache.Get([]byte(fmt.Sprintf("%d", userID)))
	if err == freecache.ErrNotFound {
		return "", app.NewError(err, app.ENotFound)
	} else if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c cacheRepository) DeleteRefreshToken(userID int64) bool {
	return c.cache.Del([]byte(fmt.Sprintf("%d", userID)))
}
