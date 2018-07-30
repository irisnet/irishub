package p2p

import (
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/p2p"
	"time"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Number of peers.
	TmMetrics 	p2p.Metrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	tmMetrics := *p2p.PrometheusMetrics()
	return &Metrics{
		tmMetrics,
	}
}

func (m *Metrics) Start(ctx app.Context) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			result := ctx.NetInfo()
			m.TmMetrics.Peers.Set(float64(result.NPeers))
		}
	}()
}
