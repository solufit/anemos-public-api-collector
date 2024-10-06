package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	libanemos "github.com/solufit/anemos-public-api-library"
)

// WriteCache writes data to Redis.
func writeCache(logger slog.Logger) {
	logger.Info("Start writing data to Cache")

	// Generate Redis Client
	client := redisClientGenerator()

	// Get Anemos Data
	data, err := getAnemosData()
	fmt.Println(data)

	if err != nil {
		logger.Error("Failed to get data from Anemos API", slog.Any("error", err))
		return
	}

	// Write data to Redis
	var dataInterface []interface{}
	err = json.Unmarshal([]byte(data), &dataInterface)
	if err != nil {
		logger.Error("Failed to unmarshal data", slog.Any("error", err))
		return
	}

	err = libanemos.CreateCache(client, dataInterface)

	if err != nil {
		logger.Error("Failed to write data to Redis", slog.Any("error", err))
		return
	}

	logger.Info("Successfully wrote data to Cache")

}

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

func getAnemosData() (string, error) {
	// Read URL from env
	anemosURL := os.Getenv("ANEMOS_URL")
	cfAccessClientID := os.Getenv("CF_ACCESS_CLIENT_ID")
	cfAccessClientSecret := os.Getenv("CF_ACCESS_CLIENT_SECRET")
	cfAccessEnabled := false
	if os.Getenv("CF_ACCESS_ENABLED") == "true" {
		cfAccessEnabled = true

		// verify cfAccessClientID and cfAccessClientSecret are set
		if cfAccessClientID == "" || cfAccessClientSecret == "" {
			return "", fmt.Errorf("CF_ACCESS_CLIENT_ID and CF_ACCESS_CLIENT_SECRET are required when CF_ACCESS_ENABLED is set to true")
		}

	} else if os.Getenv("CF_ACCESS_ENABLED") == "false" {
		cfAccessEnabled = false
	} else {
		return "", fmt.Errorf("CF_ACCESS_ENABLED is not set (true/false)")
	}

	// Get data from anemosURL

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, anemosURL, nil)

	if err != nil {
		return "", err
	}

	if cfAccessEnabled {
		req.Header.Set("CF-Access-Client-Id", cfAccessClientID)
		req.Header.Set("CF-Access-Client-Secret", cfAccessClientSecret)
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get data from Anemos API: %s", resp.Status)
	}
	defer resp.Body.Close()

	// Read response body
	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
