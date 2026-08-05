package runner

import (
	"log"
	"time"

	"github.com/notarock/a_c_a_c/pkg/metrics"
)

const maxBackoff = 5 * time.Minute

// Supervise restarts a runner when Run() returns on a fatal disconnect.
func Supervise(r *MessageCountdownRunner) {
	channel := r.Channel()
	backoff := time.Second
	for {
		start := time.Now()
		err := r.Run()
		uptime := time.Since(start)
		log.Printf("Channel %s connection ended after %s: %v — restarting in %s",
			channel, uptime, err, backoff)
		metrics.IncReconnects(channel)
		time.Sleep(backoff)
		backoff = nextBackoff(backoff, uptime)
	}
}

func nextBackoff(backoff, uptime time.Duration) time.Duration {
	if uptime > maxBackoff {
		backoff = time.Second // long-lived connection, transient failure
	}
	if backoff *= 2; backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}
