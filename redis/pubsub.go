package redis

import "github.com/garyburd/redigo/redis"
import "errors"

// NewPubSub returns a Redis pubsub connection from a config.
func NewPubSub(cnf *Config) (*redis.PubSubConn, error) {
	conn, err := New(cnf)
	if err != nil {
		return nil, err
	}

	return &redis.PubSubConn{Conn: conn}, nil
}

// New returns a new redis connection given some config.
func New(cnf *Config) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", cnf.ServerURL)
	if err != nil {
		return nil, err
	}

	data, err := redis.String(conn.Do("ECHO", "Redis OK"))
	if (err != nil) || (data != "Redis OK") {
		return nil, errors.New("could not connect to redis")
	}

	return conn, err
}
