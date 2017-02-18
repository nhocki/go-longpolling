package redis

import "github.com/garyburd/redigo/redis"

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
	return redis.Dial(
		"tcp",
		cnf.ServerURL,
		redis.DialPassword(cnf.Auth),
	)
}
