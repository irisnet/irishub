package prometheus

import (
	"github.com/irisnet/irishub/client/context"
	cs "github.com/irisnet/irishub/tools/prometheus/consensus"
	gov "github.com/irisnet/irishub/tools/prometheus/governance"
	mempl "github.com/irisnet/irishub/tools/prometheus/mempool"
	"github.com/irisnet/irishub/tools/prometheus/p2p"
	sys "github.com/irisnet/irishub/tools/prometheus/system"
)

type MetricsProvider interface {
	Start(ctx context.CLIContext)
}

type Monitor struct {
	providers []MetricsProvider
	ctx       context.CLIContext
}

func DefaultMonitor(ctx context.CLIContext) *Monitor {
	var providers []MetricsProvider
	monitor := &Monitor{
		providers: providers,
		ctx:       ctx,
	}
	monitor.AddMetricsProvider(cs.PrometheusMetrics()).
		AddMetricsProvider(p2p.PrometheusMetrics()).
		AddMetricsProvider(mempl.PrometheusMetrics()).
		AddMetricsProvider(sys.PrometheusMetrics()).
		AddMetricsProvider(gov.PrometheusMetrics())
	return monitor
}

func (m *Monitor) AddMetricsProvider(provider MetricsProvider) *Monitor {
	m.providers = append(m.providers, provider)
	return m
}

func (m *Monitor) Start() {
	for _, provider := range m.providers {
		provider.Start(m.ctx)
	}
}
