package backoff

import (
	"context"
	"testing"
	"time"
)

func Test_Backoff(t *testing.T) {
	cfg := Configuration{
		MaxDuration: time.Second * 5, // optional
	}

	b := New(&cfg)
	b.Try(context.TODO(), func() bool { return false })
}
