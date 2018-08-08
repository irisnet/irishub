package mempool

import (
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/mempool"
	"time"
)

// Metrics contains metrics exposed by this package.
// see MetricsProvider for descriptions.
type Metrics struct {
	TmMetrics mempool.Metrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	tmMetrics := *mempool.PrometheusMetrics()
	return &Metrics{
		tmMetrics,
	}
}

func (m *Metrics) Start(rpc app.Context) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			result := rpc.NumUnconfirmedTxs()
			m.TmMetrics.Size.Set(float64(result.N))
		}
	}()

}
