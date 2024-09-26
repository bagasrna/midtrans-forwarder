package redis

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Port     string
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func BuildCacheKey(key string) string {
	appName := os.Getenv("APP_NAME")
	return fmt.Sprintf("%s:%s", appName, key)
}
