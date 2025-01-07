package cache

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6380", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})
	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s\n", pong)
	return nil
}
func Save(key string, value string) error {
	ctx := context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}