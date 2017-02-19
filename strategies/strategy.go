package strategies

import (
	"io"

	"github.com/nhocki/go-longpolling/models"
)

// Strategy is an interface for pub/sub strategies
type Strategy interface {
	Setup()
	Add(s *models.Connection, channel string) error
	Remove(uuid, channel string) error
	Publish(channel string, r io.Reader) error
	TotalSubs(channel string) int
}
