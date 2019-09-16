package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReceiver             = "receiver"
	FlagReceiverOnOtherChain = "receiver-on-other-chain"
	FlagHashLock             = "hash-lock"
	FlagAmount               = "amount"
	FlagTimeLock             = "time-lock"
	FlagTimestamp            = "timestamp"
	FlagSecret               = "secret"
)

var (
	FsCreateHTLC = flag.NewFlagSet("", flag.ContinueOnError)
	FsClaimHTLC  = flag.NewFlagSet("", flag.ContinueOnError)
	FsRefundHTLC = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateHTLC.String(FlagReceiver, "", "Bech32 encoding address to receive coins")
	FsCreateHTLC.BytesHex(FlagReceiverOnOtherChain, nil, "the receiver address on the other chain")
	FsCreateHTLC.String(FlagAmount, "", "similar to the amount in the original transfer")
	FsCreateHTLC.BytesHex(FlagHashLock, nil, "the sha256 hash generated from secret (and timestamp if provided)")
	FsCreateHTLC.Uint64(FlagTimestamp, 0, "the timestamp in seconds for generating the hash lock if provided")
	FsCreateHTLC.String(FlagTimeLock, "", "the number of blocks to wait before the asset may be returned to")

	FsClaimHTLC.BytesHex(FlagHashLock, nil, "the hash lock identifying the HTLC to be claimed")
	FsClaimHTLC.String(FlagSecret, "", "the secret for generating the hash lock")

	FsRefundHTLC.BytesHex(FlagHashLock, nil, "the hash lock identifying the HTLC to be refunded")
}
