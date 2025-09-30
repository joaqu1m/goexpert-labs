package ratelimiter

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/config"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/storage"
)

type RateLimiter struct {
	storage storage.Storage
	config  *config.Config
}

type LimitResult struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetTime time.Time
}

func New(storage storage.Storage, config *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) CheckLimit(ctx context.Context, ip, token string) (*LimitResult, error) {
	if token != "" {
		return rl.checkTokenLimit(ctx, token)
	}

	return rl.checkIPLimit(ctx, ip)
}

func (rl *RateLimiter) checkTokenLimit(ctx context.Context, token string) (*LimitResult, error) {
	key := fmt.Sprintf("token:%s", token)
	blockKey := fmt.Sprintf("block:token:%s", token)

	blocked, err := rl.storage.IsBlocked(ctx, blockKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if token is blocked: %w", err)
	}

	if blocked {
		return &LimitResult{
			Allowed:   false,
			Limit:     rl.config.RateLimit.Token.RequestsPerSecond,
			Remaining: 0,
			ResetTime: time.Now().Add(rl.config.RateLimit.Token.BlockDuration),
		}, nil
	}

	timeWindow := time.Now().Unix()
	windowKey := fmt.Sprintf("%s:%d", key, timeWindow)

	count, err := rl.storage.Increment(ctx, windowKey, time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	limit := rl.config.RateLimit.Token.RequestsPerSecond
	remaining := limit - int(count)

	if count > int64(limit) {
		err = rl.storage.Block(ctx, blockKey, rl.config.RateLimit.Token.BlockDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to block token: %w", err)
		}

		return &LimitResult{
			Allowed:   false,
			Limit:     limit,
			Remaining: 0,
			ResetTime: time.Now().Add(rl.config.RateLimit.Token.BlockDuration),
		}, nil
	}

	if remaining < 0 {
		remaining = 0
	}

	return &LimitResult{
		Allowed:   true,
		Limit:     limit,
		Remaining: remaining,
		ResetTime: time.Unix(timeWindow+1, 0),
	}, nil
}

func (rl *RateLimiter) checkIPLimit(ctx context.Context, ip string) (*LimitResult, error) {
	normalizedIP := normalizeIP(ip)
	key := fmt.Sprintf("ip:%s", normalizedIP)
	blockKey := fmt.Sprintf("block:ip:%s", normalizedIP)

	blocked, err := rl.storage.IsBlocked(ctx, blockKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if IP is blocked: %w", err)
	}

	if blocked {
		return &LimitResult{
			Allowed:   false,
			Limit:     rl.config.RateLimit.IP.RequestsPerSecond,
			Remaining: 0,
			ResetTime: time.Now().Add(rl.config.RateLimit.IP.BlockDuration),
		}, nil
	}

	timeWindow := time.Now().Unix()
	windowKey := fmt.Sprintf("%s:%d", key, timeWindow)

	count, err := rl.storage.Increment(ctx, windowKey, time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	limit := rl.config.RateLimit.IP.RequestsPerSecond
	remaining := limit - int(count)

	if count > int64(limit) {
		err = rl.storage.Block(ctx, blockKey, rl.config.RateLimit.IP.BlockDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to block IP: %w", err)
		}

		return &LimitResult{
			Allowed:   false,
			Limit:     limit,
			Remaining: 0,
			ResetTime: time.Now().Add(rl.config.RateLimit.IP.BlockDuration),
		}, nil
	}

	if remaining < 0 {
		remaining = 0
	}

	return &LimitResult{
		Allowed:   true,
		Limit:     limit,
		Remaining: remaining,
		ResetTime: time.Unix(timeWindow+1, 0),
	}, nil
}

func normalizeIP(ip string) string {
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		host = ip
	}

	parsedIP := net.ParseIP(host)
	if parsedIP != nil {
		return parsedIP.String()
	}

	return host
}
