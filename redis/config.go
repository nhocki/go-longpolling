package redis

import (
	"os"
	"time"

	"github.com/nhocki/go-longpolling/models"
)

var (
	err error
)

const (
	timeout           time.Duration = 4 * time.Minute
	connectTimeout    time.Duration = 10 * time.Second
	defaultRedisConns               = 4
)

// Config holds the redis configuration
type Config struct {
	ServerURL, Auth string
	Log             models.Logger
}

// ConfigFromEnv returns a config objects that reads everything from the
// environment and adds some sensible defaults.
func ConfigFromEnv(l models.Logger) *Config {
	if l == nil {
		l = models.StdLogger
	}

	config := &Config{Log: l}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "127.0.0.1:6379"
	}
	config.ServerURL = redisURL
	return config
}
