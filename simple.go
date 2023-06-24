package backoff

import (
	"context"
	"errors"
	"time"
)

var ErrRetryLimit = errors.New("retry limit reached")

// Do attempts to call function, and then proceeds to call it on an exponential backoff until:
// - the context has been canceled
// - the retry limit has been reached
// - the function call succeeds
func Do(ctx context.Context, fn func() bool, retries int) error {
	t := time.NewTicker(1)
	defer t.Stop()

	attempt := 0

	for {
		select {
		case <-t.C:
			if fn() {
				return nil
			}

			if attempt >= retries {
				return ErrRetryLimit
			}

			t.Reset(time.Second * time.Duration(1<<attempt))
			attempt++
			continue
		case <-ctx.Done():
			return nil
		}
	}
}
