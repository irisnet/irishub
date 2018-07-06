package system

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
	"fmt"
	"os/exec"
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/disk"
)
type Metrics struct{
	CPUUtilization metrics.Gauge
	MemoUtilization metrics.Gauge
	OpenedFilesNum metrics.Gauge
	DirSize metrics.Gauge
}


// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		CPUUtilization:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "CPU_Percent",
			Help:      "CPU Utilization Percantage",
		}, []string{}),
		MemoUtilization:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "Memo_Percent",
			Help:      "Memo Utilization Percantage",
		}, []string{}),
		OpenedFilesNum:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "Opened_Files_Number",
			Help:      "Number of Opened Files, socket and other IO is included",
		}, []string{}),
		DirSize: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "system",
			Name:      "Directory_Size",
			Help:      "total size of files in given directory (in bytes)",
		}, []string{}),
	}
}

func Monitor(command string, dir_path string, metrics *Metrics)(error){
	pid, err := getPid(command)
	if err != nil{
		return err
	}
	go func(){
		for{
			time.Sleep(1*time.Second)
			metrics.RecordMetrics(int32(pid), dir_path)
		}
	}()
	return nil
}

func (metrics Metrics)RecordMetrics(pid int32, dir_path string)  {
	proc := process.Process{Pid:pid}

	if cpu_util, err := proc.CPUPercent();err != nil{
		metrics.CPUUtilization.Set(float64(-1))
	}else{
		metrics.CPUUtilization.Set(cpu_util)
	}

	if memo_util, err := proc.MemoryPercent();err != nil{
		metrics.MemoUtilization.Set(float64(-1))
	}else {
		metrics.MemoUtilization.Set(float64(memo_util))
	}

	if files, err := proc.OpenFiles();err != nil{
		metrics.OpenedFilesNum.Set(float64(-1))
	}else {
		metrics.OpenedFilesNum.Set(float64(len(files)))
	}

	if usage, err := disk.Usage(dir_path); err != nil{
		metrics.DirSize.Set(float64(-1))
	}else{
		metrics.DirSize.Set(float64(usage.Used))
	}
}

//get the pid of process that start by the given command
//the first pid return by "ps -aux|grep <command>",
// the process whose command contains "grep" is omitted
func getPid(command string)(pid int, err error){
	command_str := fmt.Sprintf("ps -aux|grep '%s'", command)
	cmd := exec.Command("/bin/bash", "-c", command_str)

	stdout, err := cmd.StdoutPipe()

	if err != nil{
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
	for _, item := range (strings.Split(string(bytes), "\n")){
		if !strings.Contains(item, "grep"){
			for j, s := range (strings.Split(item, " ")){
				if j > 0 && s != ""{
					pid, err = strconv.Atoi(s)
					if err == nil{
						return pid, nil
					}else {
						return 0, err
					}
				}
			}
		}
	}
	return 0, errors.New("cannot find the process")
}