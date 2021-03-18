// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTo                   = "to"
	FlagReceiverOnOtherChain = "receiver-on-other-chain"
	FlagSenderOnOtherChain   = "sender-on-other-chain"
	FlagAmount               = "amount"
	FlagHashLock             = "hash-lock"
	FlagTimeLock             = "time-lock"
	FlagTimestamp            = "timestamp"
	FlagSecret               = "secret"
	FlagTransfer             = "transfer"
)

var (
	FsCreateHTLC = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateHTLC.String(FlagTo, "", "Bech32 encoding address to receive tokens")
	FsCreateHTLC.String(FlagReceiverOnOtherChain, "", "Receiver address on the other chain")
	FsCreateHTLC.String(FlagSenderOnOtherChain, "", "Sender address on the other chain")
	FsCreateHTLC.String(FlagAmount, "", "Amount to be transferred")
	FsCreateHTLC.BytesHex(FlagSecret, nil, "The secret for generating the hash lock, randomly generated if omitted")
	FsCreateHTLC.BytesHex(FlagHashLock, nil, "The sha256 hash generated from secret (and timestamp if provided), generated according to the secret flag if omitted")
	FsCreateHTLC.Uint64(FlagTimestamp, 0, "The timestamp in seconds for generating the hash lock if provided")
	FsCreateHTLC.Uint64(FlagTimeLock, 0, "The number of blocks to wait before tokens may be refunded")
	FsCreateHTLC.Bool(FlagTransfer, false, "Whether it is an HTLT transaction")
}
