package strategies

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/nhocki/go-longpolling/models"
	"github.com/nhocki/go-longpolling/redis"

	client "github.com/garyburd/redigo/redis"
)

type Redis struct {
	*client.PubSubConn

	b   Strategy
	cnf *redis.Config
}

func (r *Redis) conn() (client.Conn, error) {
	return redis.New(r.cnf)
}

// NewRedisStrategy returns a long polling strategy based on Redis Pub/Sub
func NewRedisStrategy(cnf *redis.Config) (*Redis, error) {
	conn, err := redis.NewPubSub(cnf)
	if err != nil {
		return nil, err
	}

	return &Redis{
		b:          NewStdBasic(),
		cnf:        cnf,
		PubSubConn: conn,
	}, nil
}

// Setup starts the `Receive` method
func (r *Redis) Setup() {
	go r.receive()
}

// Add subscribes a subscription to it's channel
func (r *Redis) Add(c *models.Connection, channel string) error {
	if err := r.b.Add(c, channel); err != nil {
		return err
	}
	return r.Subscribe(channel)
}

// Remove removes a connection from a channel
func (r *Redis) Remove(uuid, channel string) error {
	if err := r.b.Remove(uuid, channel); err != nil {
		return err
	}

	if r.TotalSubs(channel) == 0 {
		return r.Unsubscribe(channel)
	}

	return nil
}

// Publish sends the messages to the users
func (r *Redis) Publish(channel string, rc io.Reader) error {
	c, err := r.conn()
	if err != nil {
		return err
	}
	str := readAll(rc)
	r.cnf.Log.Printf("[redis] Publishing to %s: %s", channel, str)
	c.Do("PUBLISH", channel, str)
	return c.Close()
}

// TotalSubs returns the number of subscribers in a channel
func (r *Redis) TotalSubs(channel string) int {
	return r.b.TotalSubs(channel)
}

func (r *Redis) receive() {
	for {
		switch v := r.PubSubConn.Receive().(type) {
		case client.Message:
			r.cnf.Log.Printf("[redis] Received on channel %s: %s", v.Channel, string(v.Data))
			r.b.Publish(v.Channel, bytes.NewReader(v.Data))
		case client.Subscription:
			// Do nothing
		case error:
			r.cnf.Log.Printf("error pub/sub on connection, delivery has stopped\n")
			return
		}
	}
}

func readAll(r io.Reader) string {
	b, _ := ioutil.ReadAll(r)
	return string(b)
}
