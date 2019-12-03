package types

import (
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	amt = sdk.NewInt(100)

	senderPk    = ed25519.GenPrivKey().PubKey()
	recipientPk = ed25519.GenPrivKey().PubKey()
	sender      = sdk.AccAddress(senderPk.Address())
	recipient   = sdk.AccAddress(recipientPk.Address())

	denom0   = "atom-min"
	denom1   = "btc-min"
	unidenom = FormatUniABSPrefix + "btc-min"

	input             = sdk.NewCoin(denom0, sdk.NewInt(1000))
	output            = sdk.NewCoin(denom1, sdk.NewInt(500))
	withdrawLiquidity = sdk.NewCoin(unidenom, sdk.NewInt(500))
	deadline          = time.Now().Unix()

	emptyAddr sdk.AccAddress
	emptyTime int64
)
