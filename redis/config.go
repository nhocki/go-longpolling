package redis

import (
	"net/url"
	"os"
	"time"
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
}

// ConfigFromEnv returns a config objects that reads everything from the
// environment and adds some sensible defaults.
func ConfigFromEnv() *Config {
	config := &Config{}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "127.0.0.1:6379"
	}

	srvURL, err := url.Parse(redisURL)
	if err != nil {
		panic("[ERROR] Could not parse REDIS_URL: " + redisURL)
	}

	if srvURL.User == nil {
		config.ServerURL = srvURL.String()
	} else {
		config.ServerURL = srvURL.Host
		config.Auth, _ = srvURL.User.Password()
	}

	return config
}
