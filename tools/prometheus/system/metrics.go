package system

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
)

type Metrics struct {
	CPUUtilization  metrics.Gauge
	MemoUtilization metrics.Gauge

	ProcCPUUtilization  []metrics.Gauge
	ProcMemoUtilization []metrics.Gauge
	ProcOpenedFilesNum  []metrics.Gauge
	processes       []process.Process

	DirectorySize []metrics.Gauge
	dirPaths        []string
}

func (metrics *Metrics) AddDirectory(path string) {
	metrics.dirPaths = append(metrics.dirPaths, path)
	name := fmt.Sprintf("Direcotry_Size_%s", strings.Replace(path, "/", "_", -1))
	help := fmt.Sprintf("total Size of files in %s", path)
	metrics.DirectorySize = append(metrics.DirectorySize, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      name,
		Help:      help,
	}, []string{}))
}

func (metrics *Metrics) AddProcess(command string) {
	pid, err := getPid(command)
	if err != nil {
		return
	}
	process := process.Process{Pid: int32(pid)}
	metrics.processes = append(metrics.processes, process)

	metrics.ProcCPUUtilization = append(metrics.ProcCPUUtilization, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("CPU_Percent_%d", pid),
		Help:      fmt.Sprintf("CPU Utilization Percantage of processes with pid %d", pid),
	}, []string{}))

	metrics.ProcMemoUtilization = append(metrics.ProcMemoUtilization, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("Memo_Percent_%d", pid),
		Help:      fmt.Sprintf("Memory Utilization Percantage of processes with pid %d", pid),
	}, []string{}))

	metrics.ProcOpenedFilesNum = append(metrics.ProcOpenedFilesNum, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("Opened_Files_Number_%d", pid),
		Help:      fmt.Sprintf("Number of Opened Files of processes with pid %d", pid),
	}, []string{}))

}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		CPUUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "CPU_Percent",
			Help:      "CPU Utilization Percantage",
		}, []string{}),
		MemoUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "Memo_Percent",
			Help:      "Memo Utilization Percantage",
		}, []string{}),
		ProcCPUUtilization:   make([]metrics.Gauge, 0),
		ProcMemoUtilization:   make([]metrics.Gauge, 0),
		ProcOpenedFilesNum:   make([]metrics.Gauge, 0),
		dirPaths:  make([]string, 0),
		processes: make([]process.Process, 0),
	}
}

func (metrics *Metrics) Monitor() error {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			metrics.RecordMetrics()
		}
	}()
	return nil
}

func (metrics Metrics) RecordMetrics() {

	for i, process := range(metrics.processes){
		if cpu_util, err := process.CPUPercent(); err != nil {
			metrics.ProcCPUUtilization[i].Set(float64(-1))
		} else {
			metrics.ProcCPUUtilization[i].Set(cpu_util)
		}

		if memo_util, err := process.MemoryPercent(); err != nil {
			metrics.ProcMemoUtilization[i].Set(float64(-1))
		} else {
			metrics.ProcMemoUtilization[i].Set(float64(memo_util))
		}

		if files, err := process.OpenFiles(); err != nil {
			metrics.ProcOpenedFilesNum[i].Set(float64(-1))
		} else {
			metrics.ProcOpenedFilesNum[i].Set(float64(len(files)))
		}
	}

	for i, dir_path := range (metrics.dirPaths){
		if usage, err := disk.Usage(dir_path); err != nil {
			metrics.DirectorySize[i].Set(float64(-1))
		} else {
			metrics.DirectorySize[i].Set(float64(usage.Used))
		}
	}
	vMemoStat, _ := mem.VirtualMemory()
	metrics.MemoUtilization.Set(vMemoStat.UsedPercent)

	CPUUsedPercent := float64(0.0)
	percents, _ := cpu.Percent(time.Millisecond*100, false)
	for _, percent := range (percents){
		CPUUsedPercent += percent
	}
	metrics.CPUUtilization.Set(CPUUsedPercent/4)
}

//get the pid of process that start by the given command
//the first pid return by "ps -aux|grep <command>",
// the process whose command contains "grep" is omitted
func getPid(command string) (pid int, err error) {
	command_str := fmt.Sprintf("ps -aux|grep '%s'", command)
	cmd := exec.Command("/bin/bash", "-c", command_str)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:Invalid command,", err)
		return 0, err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return 0, err
	}
	for _, item := range strings.Split(string(bytes), "\n") {
		if !strings.Contains(item, "grep") {
			for j, s := range strings.Split(item, " ") {
				if j > 0 && s != "" {
					pid, err = strconv.Atoi(s)
					if err == nil {
						return pid, nil
					} else {
						return 0, err
					}
				}
			}
		}
	}
	return 0, errors.New("cannot find the process")
}
