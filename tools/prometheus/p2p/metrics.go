package p2p

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/app"
	"github.com/pelletier/go-toml"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
	"github.com/spf13/viper"
)

/*
// Metrics contains metrics exposed by this package.
type Metrics struct {
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
*/

type Metrics struct {

	// Number of peers.
	Peers metrics.Gauge
	// Number of connected persistent peers.
	ConnectedPersistentPeers metrics.Gauge
	// Number of unconnected persistent peers.
	UnonnectedPersistentPeers metrics.Gauge
	persistent_peers          map[string]string
}

func PrometheusMetrics() *Metrics {
	return &Metrics{
		Peers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "peers",
			Help:      "Number of peers.",
		}, []string{}),
		ConnectedPersistentPeers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "connected_persistent_peers",
			Help:      "Number of connected persistent peers.",
		}, []string{}),
		UnonnectedPersistentPeers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "unconnected_persistent_peers",
			Help:      "Number of unconnected persistent peers.",
		}, []string{}),
		persistent_peers: make(map[string]string),
	}
}

func (m *Metrics) Start(ctx app.Context) {
	//func (m *Metrics) Start(ctx tools.Context) {
	m.setP2PPersistentPeers(viper.GetString("home"))
	go func() {
		for {
			time.Sleep(1 * time.Second)
			result := ctx.NetInfo()
			if result != nil{
				connected := 0
				for _, peer := range result.Peers {
					if _, exist := m.persistent_peers[string(peer.ID)]; exist {
						connected += 1
					}
				}
				m.Peers.Set(float64(result.NPeers))
				m.ConnectedPersistentPeers.Set(float64(connected))
				m.UnonnectedPersistentPeers.Set(float64(len(m.persistent_peers) - connected))
			}
		}
	}()
}

//set the p2p persistent peers by given home dir of iris config file
func (m *Metrics) setP2PPersistentPeers(homeDir string) {
	if !filepath.IsAbs(homeDir) {
		absHomeDir, err := filepath.Abs(homeDir)
		if err != nil {
			log.Println("cannot find the file ", err)
			return
		}
		homeDir = absHomeDir
	}
	configFilePath := filepath.Join(homeDir, "config/config.toml")
	//fmt.Printf("configFilePath: %s\n", configFilePath)
	if data, err := ioutil.ReadFile(configFilePath); err != nil {
		log.Println("cannot open the file ", err)
		return
	} else {
		if config, err := toml.LoadBytes(data); err != nil {
			log.Println("parse config file failed: ", err)
			return
		} else {
			persistent_peers := config.Get("p2p.persistent_peers").(string)
			for _, peer := range strings.Split(persistent_peers, ",") {
				if peer != ""{
					splited := strings.Split(peer, "@")
					m.persistent_peers[splited[0]] = splited[1]
				}
			}
		}
	}
}
