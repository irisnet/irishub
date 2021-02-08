package types

import (
	"github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRecord constructs a new Record instance
func NewRecord(txHash bytes.HexBytes, contents []Content, creator sdk.AccAddress) Record {
	return Record{
		TxHash:   txHash.String(),
		Contents: contents,
		Creator:  creator.String(),
	}
}
