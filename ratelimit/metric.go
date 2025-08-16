package ratelimit

import (
	"time"

	"arvanch/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics represents prometheus metrics for models.
type Metrics struct {
	ErrCounter *prometheus.CounterVec
	Histogram  *prometheus.HistogramVec
}

const LabelRuleName = "rule_name"

// nolint:gochecknoglobals
var (
	metrics = Metrics{
		ErrCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: config.Namespace,
				Name:      "gubernator_client_err_total",
			}, []string{LabelRuleName},
		),
		Histogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: config.Namespace,
				Name:      "gubernator_client_duration_total",
			}, []string{LabelRuleName},
		),
	}
)

func (m Metrics) report(err error, startTime time.Time, ruleName string) {
	if err != nil {
		m.ErrCounter.With(prometheus.Labels{LabelRuleName: ruleName})
	}

	m.Histogram.With(prometheus.Labels{LabelRuleName: ruleName}).Observe(time.Since(startTime).Seconds())
}
