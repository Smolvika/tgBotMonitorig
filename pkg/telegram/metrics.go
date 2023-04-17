package telegram

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "tg",
	Subsystem:  "command",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"command"})

func observeRequest(d time.Duration, command string) {
	requestMetrics.WithLabelValues(command).Observe(d.Seconds())
}
