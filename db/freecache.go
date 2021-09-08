package db

import (
	"github.com/bagus2x/recovy/config"
	"github.com/coocood/freecache"
)

func Cache(cfg *config.Config) *freecache.Cache {
	return freecache.NewCache(cfg.CacheSize())
}
