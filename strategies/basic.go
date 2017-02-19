package strategies

import (
	"io"
	"io/ioutil"
	"sync"

	"github.com/nhocki/go-longpolling/models"
)

// Basic uses an array to store users
type Basic struct {
	sync.Mutex
	log  models.Logger
	subs map[string]map[string]*models.Connection
}

// NewStdBasic returns a basic strategy that logs to STDOUT
func NewStdBasic() *Basic {
	return &Basic{
		log:  models.StdLogger,
		subs: make(map[string]map[string]*models.Connection),
	}
}

// Setup does nothing
func (b *Basic) Setup() {
	// no-op
}

// Add subscribes a subscription to it's channel
func (b *Basic) Add(c *models.Connection, channel string) error {
	b.Lock()
	defer b.Unlock()

	if b.subs == nil {
		b.subs = make(map[string]map[string]*models.Connection)
	}

	if _, ok := b.subs[channel]; !ok {
		b.subs[channel] = make(map[string]*models.Connection)
	}
	b.subs[channel][c.ID] = c

	b.log.Printf("Client %s connected to %s (curr len: %d)\n", c.ID, channel, len(b.subs[channel]))
	return nil
}

// Remove removes a connection from a channel
func (b *Basic) Remove(uuid, channel string) error {
	b.Lock()
	defer b.Unlock()

	if subscriptions, found := b.subs[channel]; found {
		if connection, f := subscriptions[uuid]; f {
			close(connection.C)
			delete(subscriptions, uuid)
		}
	}

	return nil
}

// Publish sends the messages to the users
func (b *Basic) Publish(channel string, r io.Reader) error {
	msg, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	b.log.Printf("Publishing to channel %s: %d\n", channel, len(b.subs[channel]))
	for _, conn := range b.subs[channel] {
		b.log.Printf("  - Sending to connection: %s\n", conn.ID)
		conn.C <- msg
	}
	return nil
}

// TotalSubs returns the number of subscribers in a channel
func (b *Basic) TotalSubs(channel string) int {
	b.Lock()
	defer b.Unlock()
	return len(b.subs[channel])
}
