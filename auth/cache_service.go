package auth

import "github.com/bagus2x/recovy/app"

type CacheService interface {
	GetRefreshToken(userID int64) (string, error)
}

type cacheService struct {
	cacheRepo CacheRepository
}

func NewCacheService(cacheRepo CacheRepository) CacheService {
	return &cacheService{
		cacheRepo: cacheRepo,
	}
}

func (c *cacheService) GetRefreshToken(userID int64) (string, error) {
	token, err := c.cacheRepo.GetRefreshToken(userID)
	if app.ErrorCode(err) == app.ENotFound {
		return "", app.NewError(err, app.ENotFound, "Refresh token not found")
	} else if err != nil {
		return "", err
	}

	return token, err
}
