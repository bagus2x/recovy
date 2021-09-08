package auth

import (
	"testing"

	"github.com/bagus2x/recovy/config"
	"github.com/coocood/freecache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	var cache = freecache.NewCache(100 * 1024 * 1024)
	c := NewCacheRepository(cache, config.NewTest())
	err := c.SetRefreshToken(19, "ini token")
	err2 := c.SetRefreshToken(19, "ini token")
	assert.NoError(t, err)
	assert.NoError(t, err2)
	if err != nil {
		t.Skip()
	}

	// c.DeleteRefreshToken(19)
	token, err := c.GetRefreshToken(19)
	assert.NoError(t, err)
	if err != nil {
		t.Skip()
	}
	t.Log(token)
	assert.NotZero(t, len(token))
}
