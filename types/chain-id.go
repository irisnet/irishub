package types

import (
	"math/big"
	"strings"

	etherminttypes "github.com/evmos/ethermint/types"
)

var (
	EIP155ChainID = 6688
)

func BuildEthChainID(chainID string, eip155ChainID *big.Int) string {
	if etherminttypes.IsValidChainID(chainID) {
		return chainID
	}
	chains := strings.Split(chainID, "-")
	if len(chains) != 2 {
		panic("invalid chain-id: " + chainID)
	}
	return chains[0] + "_" + eip155ChainID.String() + "-" + chains[1]
}
