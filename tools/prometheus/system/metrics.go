package system

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/app"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Metrics struct {
	CPUUtilization  metrics.Gauge
	MemoUtilization metrics.Gauge

	//processes related monitor term
	ProcCPUUtilization  metrics.Gauge
	ProcMemoUtilization metrics.Gauge
	ProcOpenedFilesNum  metrics.Gauge
	processes           process.Process
	ProcNum             metrics.Gauge
	//cmd related monitor term
	cmd string
	//storage related monitor term
	DiskUsedPercentage metrics.Gauge
	DiskFreeSpace      metrics.Gauge
	disks              string

	//file related monitor term
	FileSize      metrics.Gauge
	filePaths     string
	recursively   bool //whether compute directories size recursively
	DirectorySize metrics.Gauge
	dirPath       string
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		CPUUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "cpu_percent",
			Help:      "CPU Utilization Percantage",
		}, []string{}),
		MemoUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "memo_percent",
			Help:      "Memo Utilization Percantage",
		}, []string{}),
		ProcCPUUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "process_cpu_percent",
			Help:      "CPU Utilization Percantage of the processes iris start",
		}, []string{}),
		ProcMemoUtilization: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "process_memo_percent",
			Help:      "Memory Utilization Percantage of processes iris start",
		}, []string{}),
		ProcOpenedFilesNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "opened_files_number",
			Help:      "Number of Opened Files of processes iris start",
		}, []string{}),
		DiskUsedPercentage: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "disk_used_percentage",
			Help:      "Used Percentage of disk",
		}, []string{}),
		DiskFreeSpace: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "disk_free_space",
			Help:      "Free space of disk",
		}, []string{}),
		DirectorySize: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "direcotry_size",
			Help:      "total Size of files in home direcotry",
		}, []string{}),
		ProcNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "process_number",
			Help:      "Process number of processes iris start",
		}, []string{}),
		recursively: true,
		/*
			FileSize:            make([]metrics.Gauge, 0),
			cmd:     make([]string, 0),
			filePaths:   make([]string, 0),
			dirPaths:    make([]string, 0),
			dirPath: nil,
			disks: nil,
		*/
	}
}

func (metrics *Metrics) setPath(path string) {
	if !filepath.IsAbs(path) {
		if absPath, err := filepath.Abs(path); err != nil {
			log.Println(err.Error())
		} else {
			path = absPath
		}
	}
	if fileInfo, err := os.Stat(path); err != nil {
		log.Println(err.Error())
	} else {
		if !fileInfo.IsDir() {
			log.Println("\"" + path + "\" is not a directory!")
		}
		metrics.dirPath = path
		metrics.disks = path
	}
}

func (metrics *Metrics) add() {
	metrics.setProcess("iris start")
	home_path := viper.GetString("home")
	metrics.setPath(home_path)
	recursively := viper.GetBool("recursively")
	metrics.recursively = recursively
}

func (metrics *Metrics) setProcess(command string) {

	metrics.cmd = command
	pid, err := getPid(command)
	if err != nil {
		return
	}
	process := process.Process{Pid: int32(pid)}
	_, err = process.Cmdline()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	metrics.processes = process
}

func (metrics *Metrics) Start(ctx app.Context) {
	metrics.add()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			metrics.RecordMetrics()
		}
	}()
}

func (metrics Metrics) RecordMetrics() {
	if num, err := getProcessNum(metrics.cmd); err != nil {
		metrics.ProcNum.Set(float64(-1))
	} else {
		metrics.ProcNum.Set(float64(num))
	}

	if cpuUtil, err := metrics.processes.CPUPercent(); err != nil {
		metrics.ProcCPUUtilization.Set(float64(-1))
	} else {
		metrics.ProcCPUUtilization.Set(float64(cpuUtil))
	}

	if memoUtil, err := metrics.processes.MemoryPercent(); err != nil {
		metrics.ProcMemoUtilization.Set(float64(-1))
	} else {
		metrics.ProcMemoUtilization.Set(float64(memoUtil))
	}

	if files, err := metrics.processes.OpenFiles(); err != nil {
		metrics.ProcOpenedFilesNum.Set(float64(-1))
	} else {
		metrics.ProcOpenedFilesNum.Set(float64(len(files)))
	}

	if usage, err := disk.Usage(metrics.disks); err != nil {
		metrics.DiskUsedPercentage.Set(float64(-1))
		metrics.DiskFreeSpace.Set(float64(-1))
	} else {
		metrics.DiskUsedPercentage.Set(usage.UsedPercent)
		metrics.DiskFreeSpace.Set(float64(usage.Free))
	}

	/*
		if fileInfo, err := os.Stat(metrics.filePaths); err != nil {
			metrics.FileSize.Set(float64(-1))
		} else {
			metrics.FileSize.Set(float64(fileInfo.Size()))
		}
	*/

	if size, err := getDirSize(metrics.dirPath, metrics.recursively); err != nil {
		metrics.DirectorySize.Set(float64(-1))
	} else {
		metrics.DirectorySize.Set(float64(size))
	}

	vMemoStat, _ := mem.VirtualMemory()
	metrics.MemoUtilization.Set(vMemoStat.UsedPercent)

	CPUUsedPercent := float64(0.0)
	percents, _ := cpu.Percent(time.Millisecond*100, false)
	for _, percent := range percents {
		CPUUsedPercent += percent
	}
	metrics.CPUUtilization.Set(CPUUsedPercent)
}

//-----------------help functions-------------------------------

//get the number of process that started by the given command
func getProcessNum(command string) (num int, err error) {
	commandStr := fmt.Sprintf("ps -aux|grep '%s'|grep -v 'grep'|wc -l", command)
	cmd := exec.Command("/bin/bash", "-c", commandStr)

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
	str := string(bytes)
	if len(str) == 0 {
		return 0, err
	} else {
		str = str[:len(str)-1]
	}
	num, err = strconv.Atoi(str)

	if err == nil {
		return num, nil
	} else {
		return 0, err
	}
}

//get the pid of process that started by the given command
//the first pid return by "ps -aux|grep <command>",
// the process whose command contains "grep" is omitted
func getPid(command string) (pid int, err error) {
	commandStr := fmt.Sprintf("ps -aux|grep '%s'", command)
	cmd := exec.Command("/bin/bash", "-c", commandStr)

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

//get directory size of given path
func getDirSize(path string, recursively bool) (int64, error) {
	Separator := string(os.PathSeparator)
	queue := list.New()
	queue.PushBack(path)
	size := int64(0)
	for queue.Len() > 0 {
		path := fmt.Sprint(queue.Front().Value)
		queue.Remove(queue.Front())

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return 0, err
		}
		for _, file := range files {
			file_size := file.Size()
			size += file_size
			if file.IsDir() && recursively {
				queue.PushBack(path + Separator + file.Name())
			}
		}
	}
	return size, nil
}
