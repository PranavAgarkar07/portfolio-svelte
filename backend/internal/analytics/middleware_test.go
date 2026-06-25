package analytics

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiterConcurrentSafe(t *testing.T) {
	resetRateLimiter()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				getToken("test-ip/test-path", 100, 1*time.Minute)
			}
		}(i)
	}
	wg.Wait()
}

func TestRateLimiterAllowsUpToMax(t *testing.T) {
	resetRateLimiter()
	allowed := 0
	for i := 0; i < 10; i++ {
		if getToken("test-key", 5, 1*time.Minute) {
			allowed++
		}
	}
	if allowed != 5 {
		t.Errorf("expected 5 allowed, got %d", allowed)
	}
}

func TestRateLimiterBlocksAfterExhaustion(t *testing.T) {
	resetRateLimiter()
	blocked := false
	for i := 0; i < 10; i++ {
		if !getToken("block-key", 3, 1*time.Minute) {
			blocked = true
			break
		}
	}
	if !blocked {
		t.Error("expected rate limiter to block after exhausting tokens")
	}
}

func TestRateCleanupDoesNotPanic(t *testing.T) {
	resetRateLimiter()
	getToken("cleanup-test", 5, 1*time.Minute)
	done := make(chan struct{})
	go func() {
		defer close(done)
		startRateCleanup()
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}
}

func TestInitIsIdempotent(t *testing.T) {
	resetRateLimiter()
	resetRateLimiter()
	getToken("idempotent-test", 5, 1*time.Minute)
}
