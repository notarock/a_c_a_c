package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler() http.Handler {
	return promhttp.Handler()
}

var (
	messagesRead = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "acac_messages_read_total",
			Help: "Total number of messages read since baseline",
		},
		[]string{"channel"},
	)

	channelsTracked = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "acac_channels_tracked",
			Help: "Number of channels being tracked",
		},
	)
)

func IncMessagesRead(channel string) {
	messagesRead.WithLabelValues(channel).Inc()
}

func SetMessagesRead(channel string, value float64) {
	messagesRead.WithLabelValues(channel).Add(value)
}

func SetChannelsTracked(count int) {
	channelsTracked.Set(float64(count))
}
