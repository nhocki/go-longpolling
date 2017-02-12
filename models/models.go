package models

import (
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
)

// Connection will simulate a user that connected to a drop
type Connection struct {
	ID string

	C            chan []byte
	Disconnected <-chan bool
}

// Logger is a basic logging interface
type Logger interface {
	Printf(format string, v ...interface{})
}

// NewConnection returns a new connection with a UUID.
func NewConnection() *Connection {
	return &Connection{
		ID: uuid.NewV4().String(),
		C:  make(chan []byte),
	}
}

// StdLogger logs to STDOUT
var StdLogger Logger

func init() {
	StdLogger = log.New(os.Stdout, "", log.LstdFlags)
}
