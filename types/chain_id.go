package types

import (
	"math/big"
	"strings"

	etherminttypes "github.com/evmos/ethermint/types"
)

var (
	EIP155ChainID = "6688"
)

func BuildEthChainID(chainID string) string {
	if etherminttypes.IsValidChainID(chainID) {
		return chainID
	}

	eip155ChainID, ok := new(big.Int).SetString(EIP155ChainID, 10)
	if !ok {
		panic("invalid chain-id: " + EIP155ChainID)
	}

	chains := strings.Split(chainID, "-")
	if len(chains) != 2 {
		panic("invalid chain-id: " + chainID)
	}
	return chains[0] + "_" + eip155ChainID.String() + "-" + chains[1]
}
