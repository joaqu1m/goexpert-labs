package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis     RedisConfig
	RateLimit RateLimitConfig
	Server    ServerConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RateLimitConfig struct {
	IP    RateLimitRule
	Token RateLimitRule
}

type RateLimitRule struct {
	RequestsPerSecond int
	BlockDuration     time.Duration
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RateLimit: RateLimitConfig{
			IP: RateLimitRule{
				RequestsPerSecond: getEnvAsInt("RATE_LIMIT_IP_REQUESTS_PER_SECOND", 10),
				BlockDuration:     time.Duration(getEnvAsInt("RATE_LIMIT_IP_BLOCK_TIME_MINUTES", 5)) * time.Minute,
			},
			Token: RateLimitRule{
				RequestsPerSecond: getEnvAsInt("RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND", 100),
				BlockDuration:     time.Duration(getEnvAsInt("RATE_LIMIT_TOKEN_BLOCK_TIME_MINUTES", 5)) * time.Minute,
			},
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
