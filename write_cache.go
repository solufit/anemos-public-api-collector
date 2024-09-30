package main

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisClientGenerator is a function that generates a Redis client from env.
func redisClientGenerator() *redis.Client {
	// Read Redis Config from env
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Generate Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	return client

}
