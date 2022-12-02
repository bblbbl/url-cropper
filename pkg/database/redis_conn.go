package database

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"sync"
	"urls/pkg/etc"
)

var (
	client    *redis.Client
	redisOnce sync.Once
)

func GetRedisConnection() *redis.Client {
	redisOnce.Do(func() {
		cnf := etc.GetConfig()
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cnf.Redis.Host, cnf.Redis.Port),
			Password: cnf.Redis.Password,
			DB:       cnf.Redis.DB,
		})
	})

	return client
}

func CloseRedisConnection() {
	err := client.Close()
	if err != nil {
		etc.GetLogger().Fatalf("failed to close redis connection: %e\n", err)
	}
}
