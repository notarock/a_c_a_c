package main

import (
	"testing"
	"time"
)

// nextBackoff mirrors supervise's backoff rule.
func nextBackoff(backoff, maxBackoff, uptime time.Duration) time.Duration {
	if uptime > maxBackoff {
		backoff = time.Second
	}
	if backoff *= 2; backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}

func TestBackoff(t *testing.T) {
	const maxBackoff = 5 * time.Minute
	short := time.Second

	// Grows, caps at maxBackoff.
	b := time.Second
	for range 20 {
		b = nextBackoff(b, maxBackoff, short)
		if b > maxBackoff {
			t.Fatalf("backoff %s exceeded max %s", b, maxBackoff)
		}
	}
	if b != maxBackoff {
		t.Fatalf("backoff should have reached max %s, got %s", maxBackoff, b)
	}

	// Long uptime resets backoff.
	b = nextBackoff(b, maxBackoff, maxBackoff+time.Second)
	if b != 2*time.Second {
		t.Fatalf("backoff should reset to 2s after long uptime, got %s", b)
	}
}
