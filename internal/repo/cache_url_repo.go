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
	GetUrl(short string) (string, bool)
	GetShortUrl(url string) (string, bool)
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

func (c *RedisUrlCache) GetUrl(short string) (string, bool) {
	str, err := c.conn.Get(c.ctx, short).Result()
	if err != nil {
		return "", false
	}

	return str, true
}

func (c *RedisUrlCache) GetShortUrl(url string) (string, bool) {
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
	urlCache      urlCache
	shortUrlCache shortUrlCache
}

type urlCache struct {
	data map[string]string
	mu   sync.RWMutex
}
type shortUrlCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func UrlHashCache() *HashCache {
	once.Do(func() {
		hashCache = &HashCache{
			urlCache: urlCache{
				data: make(map[string]string),
			},
			shortUrlCache: shortUrlCache{
				data: make(map[string]string),
			},
		}
	})

	return hashCache
}

func (c *HashCache) GetUrl(short string) (string, bool) {
	c.urlCache.mu.RLock()
	defer c.urlCache.mu.RUnlock()

	v, ok := c.urlCache.data[short]

	return v, ok
}

func (c *HashCache) GetShortUrl(url string) (string, bool) {
	c.shortUrlCache.mu.RLock()
	defer c.shortUrlCache.mu.RUnlock()

	v, ok := c.shortUrlCache.data[url]

	return v, ok
}

func (c *HashCache) PutUrl(url, short string) {
	defer func() {
		c.urlCache.mu.Unlock()
		c.shortUrlCache.mu.Unlock()
	}()

	c.urlCache.mu.Lock()
	c.urlCache.data[short] = url

	c.shortUrlCache.mu.Lock()
	c.shortUrlCache.data[url] = short
}
