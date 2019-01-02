package prometheus

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/app"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestMetricsCmd(t *testing.T) {
//	sdk.InitBech32Prefix()

	cdc := app.MakeCodec()
	comm := MonitorCommand(cdc)
	viper.Set("node", "tcp://0.0.0.0:26657")
	viper.Set("address", "25C2FA00D832E8BEC64E2B5CB4AD2066ADE79DB3")
	viper.Set("account-address", "faa1gg37u8xhw5vhrfmr5mkfq8r5l4wgvd36t9hypd")
	viper.Set("home", app.DefaultNodeHome)
	viper.Set("chain-id", "test")
	viper.Set("recursively", true)
	viper.Set("trust-node", true)
	viper.Set("port", 36660)
	comm.ExecuteC()
}

var CPUUsedPercent float64
var lock *sync.RWMutex

func Test(t *testing.T) {
	CPUUtilization := prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      "cpu_percent",
		Help:      "CPU Utilization Percantage",
	}, []string{})

	lock = new(sync.RWMutex)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			RecordMetrics(CPUUtilization)
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 36660),
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("got ", err)
		}
	}()

	cmn.TrapSignal(func() {
		srv.Close()
	})
}

func getCPUUsedPercent() {
	percents, _ := cpu.Percent(time.Millisecond*1000, false)
	tmp := float64(0.0)
	for _, percent := range percents {
		tmp += percent
	}
	lock.Lock()
	defer lock.Unlock()
	CPUUsedPercent = tmp
}

func RecordMetrics(CPUUtilization metrics.Gauge) {
	go getCPUUsedPercent()
	lock.RLock()
	defer lock.RUnlock()
	CPUUtilization.Set(CPUUsedPercent)
}
