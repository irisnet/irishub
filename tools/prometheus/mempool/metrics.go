package mempool

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	tools "github.com/irisnet/irishub/tools"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
)

// Metrics contains metrics exposed by this package.
// see MetricsProvider for descriptions.
type Metrics struct {
	// Size of the mempool.
	Size metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		Size: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "mempool",
			Name:      "size",
			Help:      "Size of the mempool (number of uncommitted transactions).",
		}, []string{}),
	}
}

func (m *Metrics) Monitor(rpc tools.Context) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			result := rpc.NumUnconfirmedTxs()
			m.Size.Set(float64(result.N))
		}
	}()

}
