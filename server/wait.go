package server

import (
	"fmt"
	"net/http"
	"time"
)

func WaitStarts(url string, stopCh <-chan error) error {
	const retryAttempts = 10
	// wait until the server can start accepting connections
	for i := 0; i < retryAttempts; i++ {
		select {
		case err := <-stopCh:
			return err
		default:
			_, err := http.Get(url)
			if err == nil {
				return nil
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	return fmt.Errorf("the app can't wait until the server starts accepting connections")
}
