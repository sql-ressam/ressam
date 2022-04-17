package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// WaitStarts until the server can start accepting connections.
func WaitStarts(port int, stopCh <-chan error) error {
	const retryAttempts = 10

	// wait until the server can start accepting connections
	for i := 0; i < retryAttempts; i++ {
		select {
		case err := <-stopCh:
			return err
		default:
			_, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port))
			if err == nil {
				return nil
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	return fmt.Errorf("the app can't wait until the server starts accepting connections")
}
