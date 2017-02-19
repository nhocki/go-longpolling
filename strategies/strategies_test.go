package strategies

import (
	"sync"
	"testing"
	"time"

	"github.com/nhocki/go-longpolling/models"
)

func assertEvent(t *testing.T, c *models.Connection, wg *sync.WaitGroup) {
	select {
	case <-c.C:
		// no-op
	case <-time.After(1 * time.Second):
		t.Errorf("Did not get event")
		t.Fail()
	}
	wg.Done()
}

func assertNoEvent(t *testing.T, c *models.Connection, wg *sync.WaitGroup) {
	select {
	case <-c.C:
		t.Errorf("Got event when none was expected")
		t.Fail()
	case <-time.After(300 * time.Millisecond):
		// no-op
	}
	wg.Done()
}
