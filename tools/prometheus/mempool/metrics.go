package mempool

import (
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/mempool"
	"time"
	"log"
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
			if result, err := rpc.NumUnconfirmedTxs(); err == nil{
				m.TmMetrics.Size.Set(float64(result.N))
			}else {
				log.Println(err)
			}
		}
	}()
}