package repo

import (
	"context"
	"github.com/go-redis/redis/v9"
	"sync"
	"time"
	"urls/pkg/database"
	"urls/pkg/etc"
)

type UrlCacheRepo interface {
	GetUrl(hash string) (string, bool)
	HashByUrl(url string) (string, bool)
	PutUrl(url, short string)
}

type RedisUrlCache struct {
	conn *redis.Client
	ctx  context.Context
}

func NewUrlRedisCache(ctx context.Context) *RedisUrlCache {
	return &RedisUrlCache{
		conn: database.GetRedisConnection(),
		ctx:  ctx,
	}
}

func (c *RedisUrlCache) GetUrl(hash string) (string, bool) {
	str, err := c.conn.Get(c.ctx, hash).Result()
	if err != nil {
		return "", false
	}

	return str, true
}

func (c *RedisUrlCache) HashByUrl(url string) (string, bool) {
	str, err := c.conn.Get(c.ctx, url).Result()
	if err != nil {
		return "", false
	}

	return str, true
}

func (c *RedisUrlCache) PutUrl(url, short string) {
	err := c.conn.Set(c.ctx, url, short, time.Duration(24)*time.Hour).Err()
	if err != nil {
		etc.GetLogger().Warnf("failed to add url to cache: %e", err)
	}

	err = c.conn.Set(c.ctx, short, url, time.Duration(24)*time.Hour).Err()
	if err != nil {
		etc.GetLogger().Warnf("failed to add url to cache: %e", err)
	}
}

var (
	hashCache *HashCache
	once      sync.Once
)

type HashCache struct {
	urlCache     urlCache
	hashUrlCache hashUrlCache
}

type urlCache struct {
	data map[string]string
	mu   sync.RWMutex
}
type hashUrlCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func UrlHashCache() *HashCache {
	once.Do(func() {
		hashCache = &HashCache{
			urlCache: urlCache{
				data: make(map[string]string),
			},
			hashUrlCache: hashUrlCache{
				data: make(map[string]string),
			},
		}
	})

	return hashCache
}

func (c *HashCache) GetUrl(hash string) (string, bool) {
	c.urlCache.mu.RLock()
	defer c.urlCache.mu.RUnlock()

	v, ok := c.urlCache.data[hash]

	return v, ok
}

func (c *HashCache) HashByUrl(url string) (string, bool) {
	c.hashUrlCache.mu.RLock()
	defer c.hashUrlCache.mu.RUnlock()

	v, ok := c.hashUrlCache.data[url]

	return v, ok
}

func (c *HashCache) PutUrl(hash, url string) {
	defer func() {
		c.urlCache.mu.Unlock()
		c.hashUrlCache.mu.Unlock()
	}()

	c.urlCache.mu.Lock()
	c.urlCache.data[url] = hash

	c.hashUrlCache.mu.Lock()
	c.hashUrlCache.data[hash] = url
}
