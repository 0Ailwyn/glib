package memcache // import "github.com/carltd/glib/cache/memcache"

import (
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/carltd/glib/internal"
)

type MCache struct {
	conn *memcache.Client
}

// NewMemCache create new memcache adapter.
func newMemCache(config *internal.CacheConfig) internal.Cacher {
	return &MCache{
		conn: memcache.New(config.Dsn),
	}
}

func (c *MCache) Touch(key string, timeout time.Duration) error {
	return c.conn.Touch(key, int32(timeout/time.Second))
}

// Get get value from memcache.
func (c *MCache) Get(key string) (interface{}, error) {
	item, err := c.conn.Get(key)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

// Set set value to memcache. only support string.
func (c *MCache) Put(key string, val interface{}, timeout time.Duration) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}
	item := memcache.Item{Key: key, Value: []byte(v), Expiration: int32(timeout / time.Second)}
	return c.conn.Set(&item)
}

// Delete delete value in memcache.
func (c *MCache) Delete(key string) error {
	return c.conn.Delete(key)
}

// inc increase counter.
func (c *MCache) Increment(key string) error {
	_, err := c.conn.Increment(key, 1)
	return err
}

// dec decrease counter.
func (c *MCache) Decrement(key string) error {
	_, err := c.conn.Decrement(key, 1)
	return err
}

// FlushAll clear all cached in memcache.
func (c *MCache) ClearAll() error {
	return c.conn.FlushAll()
}

func init() {
	internal.RegisterCacheDriver("memcache", newMemCache)
}
