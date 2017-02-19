package models

import (
	"log"
	"os"

	"io/ioutil"

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

// Default loggers.
var (
	StdLogger  Logger
	NullLogger Logger
)

func init() {
	StdLogger = log.New(os.Stdout, "", log.LstdFlags)
	NullLogger = log.New(ioutil.Discard, "", log.LstdFlags)
}
