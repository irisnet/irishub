package system

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
	"github.com/irisnet/irishub/tools"
)

type Metrics struct {
	CPUUtilization  metrics.Gauge
	MemoUtilization metrics.Gauge

	//processes related monitor term
	ProcCPUUtilization  []metrics.Gauge
	ProcMemoUtilization []metrics.Gauge
	ProcOpenedFilesNum  []metrics.Gauge
	processes           []process.Process

	//cmd related monitor term
	cmd     []string
	ProcNum []metrics.Gauge

	//storage related monitor term
	DiskUsedPercentage []metrics.Gauge
	DiskFreeSpace      []metrics.Gauge
	disks              []string

	//file related monitor term
	FileSize      []metrics.Gauge
	filePaths     []string
	recursively   bool //whether compute directories size recursively
	DirectorySize []metrics.Gauge
	dirPaths      []string
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
		ProcCPUUtilization:  make([]metrics.Gauge, 0),
		ProcMemoUtilization: make([]metrics.Gauge, 0),
		ProcOpenedFilesNum:  make([]metrics.Gauge, 0),
		DiskUsedPercentage:  make([]metrics.Gauge, 0),
		DiskFreeSpace:       make([]metrics.Gauge, 0),
		FileSize:            make([]metrics.Gauge, 0),
		DirectorySize:       make([]metrics.Gauge, 0),
		processes:           make([]process.Process, 0),

		ProcNum: make([]metrics.Gauge, 0),
		cmd:     make([]string, 0),

		recursively: false,
		disks:       make([]string, 0),
		filePaths:   make([]string, 0),
		dirPaths:    make([]string, 0),
	}
}

func (metrics *Metrics) addDisk(diskPath string) {
	metrics.disks = append(metrics.disks, diskPath)

	name := fmt.Sprintf("disk_used_percentage_%s", getPathName(diskPath))
	help := fmt.Sprintf("Used Percentage of disk mount on path: %s", diskPath)
	metrics.DiskUsedPercentage = append(metrics.DiskUsedPercentage, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      name,
		Help:      help,
	}, []string{}))

	name = fmt.Sprintf("disk_free_space_%s", getPathName(diskPath))
	help = fmt.Sprintf("Free space of disk mount on path: %s", diskPath)
	metrics.DiskFreeSpace = append(metrics.DiskFreeSpace, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      name,
		Help:      help,
	}, []string{}))

}

func (metrics *Metrics) addPath(path string) {
	if fileInfo, err := os.Stat(path); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		if fileInfo.IsDir() {
			metrics.dirPaths = append(metrics.dirPaths, path)
			name := fmt.Sprintf("direcotry_size_%s", getPathName(path))
			help := fmt.Sprintf("total Size of files in %s", path)
			metrics.DirectorySize = append(metrics.DirectorySize, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
				Subsystem: "system",
				Name:      name,
				Help:      help,
			}, []string{}))
		} else {
			metrics.filePaths = append(metrics.filePaths, path)
			name := fmt.Sprintf("file_size_%s", getPathName(path))
			help := fmt.Sprintf("size of files: %s", path)
			metrics.FileSize = append(metrics.FileSize, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
				Subsystem: "system",
				Name:      name,
				Help:      help,
			}, []string{}))
		}
	}
}

func (metrics *Metrics) add() {
	commands := viper.GetString("commands")
	for _, command := range strings.Split(commands, ";") {
		if strings.TrimSpace(command) != "" {
			metrics.addProcess(strings.TrimSpace(command))
		}
	}

	disks := viper.GetString("disks")
	for _, diskPath := range strings.Split(disks, ";") {
		if strings.TrimSpace(diskPath) != "" {
			metrics.addDisk(strings.TrimSpace(diskPath))
		}
	}

	paths := viper.GetString("paths")
	for _, path := range strings.Split(paths, ";") {
		if strings.TrimSpace(path) != "" {
			metrics.addPath(strings.TrimSpace(path))
		}
	}

	recursively := viper.GetBool("recursively")
	metrics.SetRecursively(recursively)

}

func (metrics *Metrics) addProcess(command string) {
	name_command := getPathName(command)

	metrics.cmd = append(metrics.cmd, command)

	metrics.ProcNum = append(metrics.ProcNum, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("process_number_%s", name_command),
		Help:      fmt.Sprintf("Process number of processes started by command %s", command),
	}, []string{}))

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
		Name:      fmt.Sprintf("cpu_percent_%s", name_command),
		Help:      fmt.Sprintf("CPU Utilization Percantage of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

	metrics.ProcMemoUtilization = append(metrics.ProcMemoUtilization, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("memo_percent_%s", name_command),
		Help:      fmt.Sprintf("Memory Utilization Percantage of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

	metrics.ProcOpenedFilesNum = append(metrics.ProcOpenedFilesNum, prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Subsystem: "system",
		Name:      fmt.Sprintf("opened_files_number_%s", name_command),
		Help:      fmt.Sprintf("Number of Opened Files of processes with pid %d, started by command %s", pid, cmd),
	}, []string{}))

}

func (metrics *Metrics) SetRecursively(recursively bool) {
	metrics.recursively = recursively
}

func (metrics *Metrics) Start(ctx tools.Context) {
	metrics.add()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			metrics.RecordMetrics()
		}
	}()
}

func (metrics Metrics) RecordMetrics() {

	for i, cmd := range metrics.cmd {
		if num, err := getProcessNum(cmd); err != nil {
			metrics.ProcNum[i].Set(float64(-1))
		} else {
			metrics.ProcNum[i].Set(float64(num))
		}
	}

	for i, process := range metrics.processes {
		if cpuUtil, err := process.CPUPercent(); err != nil {
			metrics.ProcCPUUtilization[i].Set(float64(-1))
		} else {
			metrics.ProcCPUUtilization[i].Set(cpuUtil)
		}

		if memoUtil, err := process.MemoryPercent(); err != nil {
			metrics.ProcMemoUtilization[i].Set(float64(-1))
		} else {
			metrics.ProcMemoUtilization[i].Set(float64(memoUtil))
		}

		if files, err := process.OpenFiles(); err != nil {
			metrics.ProcOpenedFilesNum[i].Set(float64(-1))
		} else {
			metrics.ProcOpenedFilesNum[i].Set(float64(len(files)))
		}
	}

	for i, diskPath := range metrics.disks {
		if usage, err := disk.Usage(diskPath); err != nil {
			metrics.DiskUsedPercentage[i].Set(float64(-1))
			metrics.DiskFreeSpace[i].Set(float64(-1))
		} else {
			metrics.DiskUsedPercentage[i].Set(usage.UsedPercent)
			metrics.DiskFreeSpace[i].Set(float64(usage.Free))
		}
	}

	for i, filePath := range metrics.filePaths {
		if fileInfo, err := os.Stat(filePath); err != nil {
			metrics.FileSize[i].Set(float64(-1))
		} else {
			metrics.FileSize[i].Set(float64(fileInfo.Size()))
		}
	}

	for i, dirPath := range metrics.dirPaths {
		if size, err := getDirSize(dirPath, metrics.recursively); err != nil {
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
	metrics.CPUUtilization.Set(CPUUsedPercent / float64(len(percents)))
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
	if len(str) == 0{
		return 0, err
	}else{
		str = str[:len(str) - 1]
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

//conver a path to a valid Gauge monitor term name
func getPathName(path string) string {
	var buffer bytes.Buffer

	for i, ch := range path {
		if i == 0 && unicode.IsDigit(ch) {
			buffer.WriteString("_")
		}
		if unicode.IsDigit(ch) || unicode.IsLetter(ch) {
			buffer.WriteByte(byte(ch))
		} else {
			buffer.WriteByte(byte('_'))
		}
	}
	return buffer.String()
}
