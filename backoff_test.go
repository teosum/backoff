package backoff

import (
	"context"
	"testing"
	"time"
)

func Test_Backoff(t *testing.T) {
	cfg := Configuration{
		Bit:        1,                // default
		Max:        time.Second * 10, // default
		MaxRetries: 20,               // optional
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	b := New(&cfg)
	b.Try(ctx, func() bool { return false })
}
