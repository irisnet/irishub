// nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionCreateHTLC = []byte("create_htlc")

	TagAction               = sdk.TagAction
	TagSender               = "sender"
	TagReceiver             = "receiver"
	TagReceiverOnOtherChain = "receiver-on-other-chain"
	TagOutAmount            = "out-amount"
	TagInAmount             = "in-amount"
	TagSecretHashLock       = "secret-hash-lock"
	TagTimestamp            = "timestamp"
	TagExpireHeight         = "expire-height"
)
