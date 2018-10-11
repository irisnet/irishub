package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagServiceName        = "name"
	FlagServiceDescription = "service-description"
	FlagTags               = "tags"
	FlagAuthorDescription  = "author-description"
	FlagIdlContent         = "idl-content"
	FlagBroadcast          = "broadcast"
	FlagFile               = "file"
)

var (
	FsServiceName        = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceDescription = flag.NewFlagSet("", flag.ContinueOnError)
	FsTags               = flag.NewFlagSet("", flag.ContinueOnError)
	FsAuthorDescription  = flag.NewFlagSet("", flag.ContinueOnError)
	FsIdlContent         = flag.NewFlagSet("", flag.ContinueOnError)
	FsBroadcast          = flag.NewFlagSet("", flag.ContinueOnError)
	FsFile               = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsServiceName.String(FlagServiceName, "", "service name")
	FsServiceDescription.String(FlagServiceDescription, "", "service description")
	FsTags.String(FlagTags, "", "service tags")
	FsAuthorDescription.String(FlagAuthorDescription, "", "service author description")
	FsIdlContent.String(FlagIdlContent, "", "service idl content")
	FsBroadcast.String(FlagBroadcast, "", "service broadcast type")
	FsFile.String(FlagFile, "", "service idl file path")
}
