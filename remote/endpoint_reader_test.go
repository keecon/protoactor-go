package remote

import (
	"sync"
	"testing"
)

func TestEndpointReaderConnectionLifecycleIsIdempotent(t *testing.T) {
	connection := newEndpointReaderConnection()
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			connection.requestDisconnect()
		}()
		go func() {
			defer wg.Done()
			connection.finish()
		}()
	}
	wg.Wait()

	select {
	case <-connection.disconnect:
	default:
		t.Fatal("disconnect signal was not closed")
	}
	select {
	case <-connection.done:
	default:
		t.Fatal("done signal was not closed")
	}
}
