// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagURI  = "uri"
	FlagMeta = "meta"
)

// common flag sets to add to various functions
var (
	FsCreateRecord = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateRecord.String(FlagURI, "", "Source URI of the record, such as an IPFS link")
	FsCreateRecord.String(FlagMeta, "", "Metadata of the record")
}
