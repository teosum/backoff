package backoff

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Do(t *testing.T) {
	// retry cancelation handled in context
	timeout := time.Second * 10

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	// function to attempt returns true on success, and false if failed and needs to be retried
	fn := func() bool {
		fmt.Println("try")
		return false
	}

	// max number of retries - backoff exits on the first case of success, max retries, or timeout.
	retries := 3
	now := time.Now()
	_ = Do(ctx, fn, retries)

	fmt.Printf("took %v", time.Since(now))
}
