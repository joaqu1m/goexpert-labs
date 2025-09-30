package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/config"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/ratelimiter"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) Get(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) Reset(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRateLimiter_CheckIPLimit_Allowed(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IP: config.RateLimitRule{
				RequestsPerSecond: 10,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	rateLimiter := ratelimiter.New(mockStorage, cfg)

	mockStorage.On("IsBlocked", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	mockStorage.On("Increment", mock.Anything, mock.AnythingOfType("string"), time.Second).Return(int64(5), nil)

	ctx := context.Background()
	result, err := rateLimiter.CheckLimit(ctx, "192.168.1.1", "")

	assert.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 5, result.Remaining)
	mockStorage.AssertExpectations(t)
}

func TestRateLimiter_CheckIPLimit_Exceeded(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IP: config.RateLimitRule{
				RequestsPerSecond: 10,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	rateLimiter := ratelimiter.New(mockStorage, cfg)

	mockStorage.On("IsBlocked", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	mockStorage.On("Increment", mock.Anything, mock.AnythingOfType("string"), time.Second).Return(int64(11), nil)
	mockStorage.On("Block", mock.Anything, mock.AnythingOfType("string"), 5*time.Minute).Return(nil)

	ctx := context.Background()
	result, err := rateLimiter.CheckLimit(ctx, "192.168.1.1", "")

	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Remaining)
	mockStorage.AssertExpectations(t)
}

func TestRateLimiter_CheckTokenLimit_Allowed(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			Token: config.RateLimitRule{
				RequestsPerSecond: 100,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	rateLimiter := ratelimiter.New(mockStorage, cfg)

	mockStorage.On("IsBlocked", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	mockStorage.On("Increment", mock.Anything, mock.AnythingOfType("string"), time.Second).Return(int64(50), nil)

	ctx := context.Background()
	result, err := rateLimiter.CheckLimit(ctx, "192.168.1.1", "token123")

	assert.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 100, result.Limit)
	assert.Equal(t, 50, result.Remaining) // 100 - 50 = 50
	mockStorage.AssertExpectations(t)
}

func TestRateLimiter_CheckTokenLimit_AlreadyBlocked(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			Token: config.RateLimitRule{
				RequestsPerSecond: 100,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	rateLimiter := ratelimiter.New(mockStorage, cfg)

	mockStorage.On("IsBlocked", mock.Anything, mock.AnythingOfType("string")).Return(true, nil)

	ctx := context.Background()
	result, err := rateLimiter.CheckLimit(ctx, "192.168.1.1", "token123")

	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 100, result.Limit)
	assert.Equal(t, 0, result.Remaining)
	mockStorage.AssertExpectations(t)
}

func TestRateLimiter_TokenPrecedenceOverIP(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IP: config.RateLimitRule{
				RequestsPerSecond: 10,
				BlockDuration:     5 * time.Minute,
			},
			Token: config.RateLimitRule{
				RequestsPerSecond: 100,
				BlockDuration:     5 * time.Minute,
			},
		},
	}

	rateLimiter := ratelimiter.New(mockStorage, cfg)

	mockStorage.On("IsBlocked", mock.Anything, mock.MatchedBy(func(key string) bool {
		return key == "block:token:token123"
	})).Return(false, nil)
	mockStorage.On("Increment", mock.Anything, mock.MatchedBy(func(key string) bool {
		return key != ""
	}), time.Second).Return(int64(50), nil)

	ctx := context.Background()
	result, err := rateLimiter.CheckLimit(ctx, "192.168.1.1", "token123")

	assert.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 100, result.Limit)
	assert.Equal(t, 50, result.Remaining)
	mockStorage.AssertExpectations(t)
}
