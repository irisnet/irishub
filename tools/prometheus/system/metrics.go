package system

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"container/list"
	"bytes"
	"unicode"
)

type Metrics struct {
	CPUUtilization  metrics.Gauge
	MemoUtilization metrics.Gauge

	//processes related monitor term
	ProcCPUUtilization  []metrics.Gauge
	ProcMemoUtilization []metrics.Gauge
	ProcOpenedFilesNum  []metrics.Gauge
	processes           []process.Process

	//storage related monitor term
	DiskUsedPercentage []metrics.Gauge
	DiskFreeSpace      []metrics.Gauge
	disks              []string

	//file related monitor term
	FileSize      []metrics.Gauge
	filePaths     []string
	recursively   bool  //whether compute directories size recursively
	DirectorySize []metrics.Gauge
	dirPaths      []string
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
		ProcCPUUtilization:  make([]metrics.Gauge, 0),
		ProcMemoUtilization: make([]metrics.Gauge, 0),
		ProcOpenedFilesNum:  make([]metrics.Gauge, 0),
		DiskUsedPercentage:  make([]metrics.Gauge, 0),
		DiskFreeSpace:       make([]metrics.Gauge, 0),
		FileSize:            make([]metrics.Gauge, 0),
		DirectorySize:       make([]metrics.Gauge, 0),
		processes:           make([]process.Process, 0),
		recursively:         false,
		disks:               make([]string, 0),
		filePaths:           make([]string, 0),
		dirPaths:            make([]string, 0),
	}
}

func (metrics *Metrics) AddDisk(disk_path string) {
	metrics.disks = append(metrics.disks, disk_path)

	name := fmt.Sprintf("Disk_Used_Percentage_%s", get_path_name(disk_path))
	help := fmt.Sprintf("Used Percentage of disk mount on path: %s", disk_path)
	metrics.DiskUsedPercentage = append(metrics.DiskUsedPercentage, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      name,
		Help:      help,
	}, []string{}))

	name = fmt.Sprintf("Disk_Free_Space_%s", get_path_name(disk_path))
	help = fmt.Sprintf("Free space of disk mount on path: %s", disk_path)
	metrics.DiskFreeSpace = append(metrics.DiskFreeSpace, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      name,
		Help:      help,
	}, []string{}))

}

func (metrics *Metrics) AddPath(path string) {
	if fileInfo, err := os.Stat(path); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		if fileInfo.IsDir() {
			metrics.dirPaths = append(metrics.dirPaths, path)
			name := fmt.Sprintf("Direcotry_Size_%s", get_path_name(path))
			//name := fmt.Sprintf("Direcotry_Size_%s", strings.Replace(path, "/", "_", -1))
			help := fmt.Sprintf("total Size of files in %s", path)
			metrics.DirectorySize = append(metrics.DirectorySize, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
				Subsystem: "system",
				Name:      name,
				Help:      help,
			}, []string{}))
		} else {
			metrics.filePaths = append(metrics.filePaths, path)
			name := fmt.Sprintf("File_Size_%s", get_path_name(path))
			help := fmt.Sprintf("size of files: %s", path)
			metrics.FileSize = append(metrics.FileSize, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
				Subsystem: "system",
				Name:      name,
				Help:      help,
			}, []string{}))
		}
	}
}

func (metrics *Metrics) AddProcess(command string) {
	pid, err := getPid(command)
	if err != nil {
		return
	}
	process := process.Process{Pid: int32(pid)}
	cmd, err := process.Cmdline()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	metrics.processes = append(metrics.processes, process)

	metrics.ProcCPUUtilization = append(metrics.ProcCPUUtilization, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("CPU_Percent_%d", pid),
		Help:      fmt.Sprintf("CPU Utilization Percantage of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

	metrics.ProcMemoUtilization = append(metrics.ProcMemoUtilization, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("Memo_Percent_%d", pid),
		Help:      fmt.Sprintf("Memory Utilization Percantage of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

	metrics.ProcOpenedFilesNum = append(metrics.ProcOpenedFilesNum, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("Opened_Files_Number_%d", pid),
		Help:      fmt.Sprintf("Number of Opened Files of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

}

func (metrics *Metrics)SetRecursively(recursively bool)  {
	metrics.recursively = recursively
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

	for i, process := range metrics.processes {
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

	for i, disk_path := range metrics.disks {
		if usage, err := disk.Usage(disk_path); err != nil {
			metrics.DirectorySize[i].Set(float64(-1))
			metrics.DiskFreeSpace[i].Set(float64(-1))
		} else {
			metrics.DirectorySize[i].Set(usage.UsedPercent)
			metrics.DiskFreeSpace[i].Set(float64(usage.Free))
		}
	}

	for i, file_path := range metrics.filePaths {
		if fileInfo, err := os.Stat(file_path); err != nil {
			metrics.FileSize[i].Set(float64(-1))
		} else {
			metrics.FileSize[i].Set(float64(fileInfo.Size()))
		}
	}

	for i, dir_path := range metrics.dirPaths {
		if size, err := get_dir_size(dir_path, metrics.recursively); err != nil {
			metrics.DirectorySize[i].Set(float64(-1))
		} else {
			metrics.DirectorySize[i].Set(float64(size))
		}
	}

	vMemoStat, _ := mem.VirtualMemory()
	metrics.MemoUtilization.Set(vMemoStat.UsedPercent)

	CPUUsedPercent := float64(0.0)
	percents, _ := cpu.Percent(time.Millisecond*100, false)
	for _, percent := range percents {
		CPUUsedPercent += percent
	}
	metrics.CPUUtilization.Set(CPUUsedPercent / 4)
}

//-----------------help functions-------------------------------

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

//get directory size of given path
func get_dir_size(dir_path string, recursively bool) (int64, error) {
	Separator := string(os.PathSeparator)
	queue := list.New()
	queue.PushBack(dir_path)
	size := int64(0)
	for ;queue.Len() > 0;{
		path := fmt.Sprint(queue.Front().Value)
		queue.Remove(queue.Front())

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return 0, err
		}
		for _, file := range files {
			file_size := file.Size()
			size += file_size
			if file.IsDir() && recursively{
				queue.PushBack(path + Separator + file.Name())
			}
		}
	}
	return size, nil
}

//conver a path to a valid Gauge monitor term name
func get_path_name(path string)(path_name string){
	var buffer bytes.Buffer

	for i, ch := range path{
		if i == 0 && unicode.IsDigit(ch){
			buffer.WriteString("_")
		}
		if unicode.IsDigit(ch) || unicode.IsLetter(ch){
			buffer.WriteByte(byte(ch))
		}else {
			buffer.WriteByte(byte('_'))
		}
	}
	return buffer.String()
}