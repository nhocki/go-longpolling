package strategies

import (
	"bytes"
	"sync"
	"testing"

	"github.com/nhocki/go-longpolling/models"
	"github.com/nhocki/go-longpolling/redis"
	"github.com/stretchr/testify/assert"
)

func newRedisStrategy(t *testing.T) *Redis {
	cnf := redis.ConfigFromEnv(models.NullLogger)
	conn, err := redis.NewPubSub(cnf)
	if err != nil {
		assert.Fail(t, "Error! %s", err)
		return nil
	}

	return &Redis{
		b:          &Basic{log: models.NullLogger},
		cnf:        cnf,
		PubSubConn: conn,
	}
}

func TestSingleRedisStrategy(t *testing.T) {
	strategy := newRedisStrategy(t)
	assert.NotNil(t, strategy)
	strategy.Setup()

	events := models.NewConnection()
	messages := models.NewConnection()

	assert.NoError(t, strategy.Add(events, "events"))
	assert.NoError(t, strategy.Add(messages, "messages"))

	assert.Equal(t, 1, strategy.TotalSubs("events"))
	assert.Equal(t, 1, strategy.TotalSubs("messages"))

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go assertEvent(t, events, wg)
	go assertNoEvent(t, messages, wg)
	strategy.Publish("events", bytes.NewReader([]byte(`Hello, World!`)))
	wg.Wait()

	wg.Add(2)
	go assertEvent(t, messages, wg)
	go assertNoEvent(t, events, wg)
	strategy.Publish("messages", bytes.NewReader([]byte(`Hello, World!`)))
	wg.Wait()

	assert.Equal(t, 1, strategy.TotalSubs("events"))
	assert.Equal(t, 1, strategy.TotalSubs("messages"))

	assert.NoError(t, strategy.Remove(events.ID, "events"))
	assert.NoError(t, strategy.Remove(messages.ID, "events"))
	assert.NoError(t, strategy.Remove(messages.ID, "messages"))
	assert.NoError(t, strategy.Remove(messages.ID, "messages"))

	assert.Equal(t, 0, strategy.TotalSubs("events"))
	assert.Equal(t, 0, strategy.TotalSubs("messages"))
}

func TestMultipleRedisStrategy(t *testing.T) {
	publisher := newRedisStrategy(t)
	assert.NotNil(t, publisher)

	subscriber := newRedisStrategy(t)
	assert.NotNil(t, subscriber)

	// Start subscriber listener
	subscriber.Setup()

	events := models.NewConnection()
	messages := models.NewConnection()

	// Add connections to subscriber
	assert.NoError(t, subscriber.Add(events, "events"))
	assert.NoError(t, subscriber.Add(messages, "messages"))

	// Send messages to publisher. Connections should still get the messages.
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go assertEvent(t, events, wg)
	go assertNoEvent(t, messages, wg)
	publisher.Publish("events", bytes.NewReader([]byte(`Hello, World!`)))
	wg.Wait()

	wg.Add(2)
	go assertEvent(t, messages, wg)
	go assertNoEvent(t, events, wg)
	publisher.Publish("messages", bytes.NewReader([]byte(`Hello, World!`)))
	wg.Wait()
}
