package storage

import (
	"context"
	"time"
)

// Storage defines the interface for storing rate limiter data
type Storage interface {
	// Increment increments the counter for the given key and returns the new value
	// If the key doesn't exist, it creates it with value 1 and sets the expiration
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)

	// Get returns the current count for the given key
	Get(ctx context.Context, key string) (int64, error)

	// Reset resets the counter for the given key
	Reset(ctx context.Context, key string) error

	// IsBlocked checks if a key is currently blocked
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Block blocks a key for the specified duration
	Block(ctx context.Context, key string, duration time.Duration) error

	// Close closes the storage connection
	Close() error
}
