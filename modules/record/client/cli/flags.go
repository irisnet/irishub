// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagURI  = "uri"
	FlagMeta = "meta"
)

// common flagsets to add to various functions
var (
	FsCreateRecord = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateRecord.String(FlagURI, "", "source uri of record, such as an ipfs link")
	FsCreateRecord.String(FlagMeta, "", "meta data of record")
}
