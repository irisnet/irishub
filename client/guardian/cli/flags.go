package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagProfilerAddress = "profiler-address"
	FlagProfilerName    = "profiler-name"
)

var (
	FsProfiler = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsProfiler.String(FlagProfilerAddress, "", "bech32 encoded account of the profiler")
	FsProfiler.String(FlagProfilerName, "", "name of the profiler")
}
