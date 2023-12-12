package types

import (
	"fmt"
	"math/big"
)

var (
	EIP155ChainID = "6688"
)

func parseChainID(_ string) (*big.Int, error) {
	eip155ChainID, ok := new(big.Int).SetString(EIP155ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain-id: %s" + EIP155ChainID)
	}
	return eip155ChainID, nil
}
