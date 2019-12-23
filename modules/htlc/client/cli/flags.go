// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReceiverOnOtherChain = "receiver-on-other-chain"
	FlagHashLock             = "hash-lock"
	FlagAmount               = "amount"
	FlagTimeLock             = "time-lock"
	FlagTimestamp            = "timestamp"
	FlagSecret               = "secret"
	FlagTo                   = "to"
)

// common flagsets to add to various functions
var (
	FsCreateHTLC = flag.NewFlagSet("", flag.ContinueOnError)
	FsClaimHTLC  = flag.NewFlagSet("", flag.ContinueOnError)
	FsRefundHTLC = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateHTLC.String(FlagTo, "", "Bech32 encoding address to receive coins")
	FsCreateHTLC.String(FlagReceiverOnOtherChain, "", "The claim receiving address on the other chain")
	FsCreateHTLC.String(FlagAmount, "", "Similar to the amount in the original transfer")
	FsCreateHTLC.BytesHex(FlagSecret, nil, "The secret for generating the hash lock, randomly generated if omitted")
	FsCreateHTLC.BytesHex(FlagHashLock, nil, "The sha256 hash generated from secret (and timestamp if provided), generated from the secret flag if omitted")
	FsCreateHTLC.Uint64(FlagTimestamp, 0, "The timestamp in seconds for generating the hash lock if provided")
	FsCreateHTLC.String(FlagTimeLock, "", "The number of blocks to wait before the asset may be returned to")
	FsClaimHTLC.BytesHex(FlagHashLock, nil, "The hash lock identifying the HTLC to be claimed")
	FsClaimHTLC.BytesHex(FlagSecret, nil, "The secret for generating the hash lock")
	FsRefundHTLC.BytesHex(FlagHashLock, nil, "The hash lock identifying the HTLC to be refunded")
}
