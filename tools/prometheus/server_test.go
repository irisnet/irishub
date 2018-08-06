package prometheus

import (
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/app"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"testing"
	"github.com/go-kit/kit/metrics"
	"time"
	"github.com/shirou/gopsutil/cpu"
	"net/http"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	cmn "github.com/tendermint/tendermint/libs/common"
	"sync"
)

func TestMetricsCmd(t *testing.T) {
	cdc := app.MakeCodec()
	comm := MonitorCommand(cdc)
	viper.Set("node", "tcp://0.0.0.0:26657")
	viper.Set("address", "92F517F241560761B5D201C0E103C1C3D5C277DB")
	viper.Set("home", app.DefaultNodeHome)
	viper.Set("chain-id", "fuxi-test")
	viper.Set("recursively", true)
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

func getCPUUsedPercent(){
	percents, _ := cpu.Percent(time.Millisecond*1000, false)
	tmp := float64(0.0)
	for _, percent := range percents {
		tmp += percent
	}
	lock.Lock()
	defer lock.Unlock()
	CPUUsedPercent = tmp
}

func RecordMetrics(CPUUtilization (metrics.Gauge)) {
	go getCPUUsedPercent()
	lock.RLock()
	defer lock.RUnlock()
	CPUUtilization.Set(CPUUsedPercent)
}

func TestCond(t *testing.T){
	lock := new(sync.Once)
}
