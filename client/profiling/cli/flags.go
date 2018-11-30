package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagProfilerAddress = "profiler-address"
	FlagProfilerName    = "profiler-name"
)

var (
	FsProfilerAddress = flag.NewFlagSet("", flag.ContinueOnError)
	FsProfilerName    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsProfilerAddress.String(FlagProfilerAddress, "", "bech32 encoded account of the profiler")
	FsProfilerName.String(FlagProfilerName, "", "name of the profiler")
}
