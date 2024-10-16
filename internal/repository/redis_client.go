package repository

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
