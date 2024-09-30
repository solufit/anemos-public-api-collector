package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisClientGenerator(t *testing.T) {
	// Mock environment variables
	mockRedisHost := "localhost"
	mockRedisPort := "6379"
	mockRedisPassword := "password"

	os.Setenv("REDIS_HOST", mockRedisHost)
	os.Setenv("REDIS_PORT", mockRedisPort)
	os.Setenv("REDIS_PASSWORD", mockRedisPassword)

	// Generate Redis client
	client := redisClientGenerator()

	// Check if the client is not nil
	assert.NotNil(t, client, "Expected Redis client to be non-nil")

	// Check if the client has the correct address
	expectedAddr := fmt.Sprintf("%s:%s", mockRedisHost, mockRedisPort)
	assert.Equal(t, expectedAddr, client.Options().Addr, "Expected Redis client address to match")

	// Check if the client has the correct password
	assert.Equal(t, mockRedisPassword, client.Options().Password, "Expected Redis client password to match")
}
