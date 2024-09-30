package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetAnemosData(t *testing.T) {
	// Mock environment variables
	mockAnemosURL := "http://mock-anemos-url"
	mockCFClientID := "mock-client-id"
	mockCFClientSecret := "mock-client-secret"

	os.Setenv("ANEMOS_URL", mockAnemosURL)
	os.Setenv("CF_ACCESS_CLIENT_ID", mockCFClientID)
	os.Setenv("CF_ACCESS_CLIENT_SECRET", mockCFClientSecret)

	t.Run("CF_ACCESS_ENABLED not set", func(t *testing.T) {
		os.Unsetenv("CF_ACCESS_ENABLED")
		_, err := getAnemosData()
		assert.Error(t, err, "Expected error when CF_ACCESS_ENABLED is not set")
	})

	t.Run("CF_ACCESS_ENABLED set to true without credentials", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "true")
		os.Unsetenv("CF_ACCESS_CLIENT_ID")
		os.Unsetenv("CF_ACCESS_CLIENT_SECRET")
		_, err := getAnemosData()
		assert.Error(t, err, "Expected error when CF_ACCESS_ENABLED is true but credentials are missing")
	})
	t.Run("CF_ACCESS_ENABLED set to true without Client ID", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "true")
		os.Unsetenv("CF_ACCESS_CLIENT_ID")
		_, err := getAnemosData()
		assert.Error(t, err, "Expected error when CF_ACCESS_ENABLED is true but credentials are missing")
	})

	t.Run("CF_ACCESS_ENABLED set to true without Client Secret", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "true")
		os.Unsetenv("CF_ACCESS_CLIENT_SECRET")
		_, err := getAnemosData()
		assert.Error(t, err, "Expected error when CF_ACCESS_ENABLED is true but credentials are missing")
	})

	t.Run("CF_ACCESS_ENABLED set to true with credentials", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "true")
		os.Setenv("CF_ACCESS_CLIENT_ID", mockCFClientID)
		os.Setenv("CF_ACCESS_CLIENT_SECRET", mockCFClientSecret)

		// Mock server to simulate Anemos API
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, mockCFClientID, r.Header.Get("CF-Access-Client-Id"))
			assert.Equal(t, mockCFClientSecret, r.Header.Get("CF-Access-Client-Secret"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("mock response"))
		}))
		defer mockServer.Close()

		os.Setenv("ANEMOS_URL", mockServer.URL)

		data, err := getAnemosData()
		assert.NoError(t, err, "Expected no error when CF_ACCESS_ENABLED is true with valid credentials")
		assert.Equal(t, "mock response", data, "Expected mock response from Anemos API")
	})

	t.Run("CF_ACCESS_ENABLED set to false", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "false")

		// Mock server to simulate Anemos API
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("mock response"))
		}))
		defer mockServer.Close()

		os.Setenv("ANEMOS_URL", mockServer.URL)

		data, err := getAnemosData()
		assert.NoError(t, err, "Expected no error when CF_ACCESS_ENABLED is false")
		assert.Equal(t, "mock response", data, "Expected mock response from Anemos API")
	})

	t.Run("Anemos API returns non-200 status", func(t *testing.T) {
		os.Setenv("CF_ACCESS_ENABLED", "false")

		// Mock server to simulate Anemos API
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		os.Setenv("ANEMOS_URL", mockServer.URL)

		_, err := getAnemosData()
		assert.Error(t, err, "Expected error when Anemos API returns non-200 status")
	})
}
