package types

import (
	"github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRecord constructs a record
func NewRecord(txHash bytes.HexBytes, contents []Content, creator sdk.AccAddress) Record {
	return Record{
		TxHash:   txHash,
		Contents: contents,
		Creator:  creator,
	}
}
