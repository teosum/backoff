package backoff

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// Backoff is an exponential backoff function runner.
type Backoff struct {
	cfg     *Configuration
	retries uint64
	maxed   bool
	d       time.Duration
}

// Configuration for Backoff.
type Configuration struct {
	// Bit is the zero base int to increase exponentially (2 ^ x).
	Bit uint

	// Max duration a single backoff can wait. A zero value will waiting indefinitely.
	Max time.Duration

	// Max duration a series of backoffs can wait before stopping the additional retries
	// and returning error. If zero, there is no duration limit.
	MaxDuration time.Duration

	// Max number of attempts the backoff can try before returning error. If zero, an
	// unlimited number of attempts can be made.
	MaxRetries uint64
}

func (c *Configuration) validate() {
	if c.Bit == 0 {
		fmt.Println("Defaulting Bit to 1")
		c.Bit = 1
	}

	if c.Max < 1 {
		fmt.Println("Defaulting Max to 10s")
		c.Max = time.Second * 10
	}
}

// New returns a pointer to a created backoff. An optional configuration object can be
// passed in. By default, the backoff with increase exponentially by a base of 2, with
// a max backoff of 10 seconds.
func New(cfg *Configuration) *Backoff {
	if cfg == nil {
		cfg = &Configuration{}
	}

	cfg.validate()
	return &Backoff{cfg: cfg}
}

var (
	ErrRetryLimit  = errors.New("max number of backoff retries reached")
	ErrRetryTimout = errors.New("backoff timeout reached")
)

// Try attempts to invoke a boolean function. An error is returned if the number of
// reties exceeds the max configurable limit.
func (b *Backoff) Try(ctx context.Context, fn func() bool) error {
	t := time.NewTimer(0)
	defer t.Stop()
	defer b.Reset()

	var tCtx context.Context
	var tCancel context.CancelFunc
	if b.cfg.MaxDuration > 0 {
		tCtx, tCancel = context.WithTimeout(ctx, b.cfg.MaxDuration)
		defer tCancel()
	}

	for {
		select {
		case <-t.C:
			if !fn() {
				if !b.maxed {
					d := time.Second * time.Duration(b.cfg.Bit<<b.retries)
					if d > b.cfg.Max {
						b.d = b.cfg.Max
						b.maxed = true
					} else {
						b.d = d
					}
				}

				atomic.AddUint64(&b.retries, 1)
				if b.cfg.MaxRetries > 0 && atomic.LoadUint64(&b.retries) >= b.cfg.MaxRetries {
					return ErrRetryLimit
				}

				fmt.Printf("Retry %d - waiting %s\n", b.retries, b.d)
				t.Reset(b.d)
				continue
			}

			return nil

		case <-tCtx.Done():
			return ErrRetryTimout
		case <-ctx.Done():
			return nil
		}
	}
}

// Reset the backoff attempt counter.
func (b *Backoff) Reset() {
	b.retries = 0
}
