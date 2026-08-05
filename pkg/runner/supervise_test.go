package runner

import (
	"testing"
	"time"
)

func TestNextBackoff(t *testing.T) {
	short := time.Second

	// Grows, caps at maxBackoff.
	b := time.Second
	for range 20 {
		b = nextBackoff(b, short)
		if b > maxBackoff {
			t.Fatalf("backoff %s exceeded max %s", b, maxBackoff)
		}
	}
	if b != maxBackoff {
		t.Fatalf("backoff should reach max %s, got %s", maxBackoff, b)
	}

	// Long uptime resets backoff.
	if b = nextBackoff(b, maxBackoff+time.Second); b != 2*time.Second {
		t.Fatalf("backoff should reset to 2s, got %s", b)
	}
}
