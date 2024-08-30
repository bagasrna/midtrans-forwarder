package config

import (
	"os"

	"midtrans-forwarder/pkg/database"
	"midtrans-forwarder/pkg/redis"
)

type Config struct {
	DBConfig         database.MySQLConfig
	RedisConfig      redis.RedisConfig
	MidtransServerKey string
}

func Load() Config {
	return Config{
		DBConfig: database.MySQLConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		RedisConfig: redis.RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
			Port: os.Getenv("REDIS_PORT"),
		},
		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
	}
}