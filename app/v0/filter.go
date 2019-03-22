package v0

import (
	abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/version"
)

// need to set filter_peers = true in the config.toml to enable the filter
// reject all peers for the testing
func FilterPeerByPubKey(info string) abci.ResponseQuery {

	return abci.ResponseQuery{
		Code:      uint32(sdk.CodeInvalidPubKey),
		Codespace: string(sdk.CodespaceRoot),
		Value:     []byte(version.GetVersion()),
	}
}
