package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/config"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/middleware"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/ratelimiter"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/storage"
)

func TestIntegration_RateLimiter_BasicFunctionality(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		Redis: config.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		},
		RateLimit: config.RateLimitConfig{
			IP: config.RateLimitRule{
				RequestsPerSecond: 5,
				BlockDuration:     1 * time.Minute,
			},
			Token: config.RateLimitRule{
				RequestsPerSecond: 10,
				BlockDuration:     1 * time.Minute,
			},
		},
	}

	redisStorage := storage.NewRedisStorage(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	ctx := context.Background()
	err := redisStorage.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer redisStorage.Close()

	rateLimiter := ratelimiter.New(redisStorage, cfg)
	middleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	t.Run("IP Rate Limiting", func(t *testing.T) {
		testIP := "192.168.1.100"

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Real-IP", testIP)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "5", w.Header().Get("X-RateLimit-Limit"))
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Real-IP", testIP)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "maximum number of requests")
	})

	t.Run("Token Rate Limiting", func(t *testing.T) {
		testToken := "test-token-123"

		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("API_KEY", testToken)
			req.Header.Set("X-Real-IP", "192.168.1.200")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", testToken)
		req.Header.Set("X-Real-IP", "192.168.1.200")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("Token Precedence Over IP", func(t *testing.T) {
		testToken := "precedence-token"
		testIP := "192.168.1.201"

		for i := 0; i < 6; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("API_KEY", testToken)
			req.Header.Set("X-Real-IP", testIP)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
		}
	})
}

func TestIntegration_FullServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Configuration Loading", func(t *testing.T) {
		cfg, err := config.Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.NotEmpty(t, cfg.Server.Port)
	})
}

func BenchmarkRateLimiter_CheckLimit(b *testing.B) {
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IP: config.RateLimitRule{
				RequestsPerSecond: 1000,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	storage := &InMemoryStorage{
		data: make(map[string]int64),
	}

	rateLimiter := ratelimiter.New(storage, cfg)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = rateLimiter.CheckLimit(ctx, "192.168.1.1", "")
		}
	})
}

type InMemoryStorage struct {
	data map[string]int64
}

func (s *InMemoryStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	s.data[key]++
	return s.data[key], nil
}

func (s *InMemoryStorage) Get(ctx context.Context, key string) (int64, error) {
	return s.data[key], nil
}

func (s *InMemoryStorage) Reset(ctx context.Context, key string) error {
	delete(s.data, key)
	return nil
}

func (s *InMemoryStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (s *InMemoryStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	return nil
}

func (s *InMemoryStorage) Close() error {
	return nil
}
