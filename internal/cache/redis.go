package cache

import (
	"context"
	"github.com/go-redis/redis/v8"

	"log"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisCache{client: client}
}

func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil // Key doesn't exist in cache
	}
	if err != nil {
		log.Println("Redis GET error:", err)
		return "", err
	}
	return val, nil
}

func (r *RedisCache) Set(key string, value string) error {
	err := r.client.Set(context.Background(), key, value, 24*time.Second).Err()
	if err != nil {
		log.Println("Redis SET error:", err)
		return err
	}
	return nil
}
