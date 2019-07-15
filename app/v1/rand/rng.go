package rand

import (
	sdk "github.com/irisnet/irishub/types"
)

// RNG is a random number generator
type RNG interface {
	GetRand() sdk.Dec // interface which returns a random number between (0,1)
}

// PRNG represents a pseudo-random number implementation based on block related data for RNG
type PRNG struct {
	BlockHash      []byte
	BlockTimestamp int64
	TxInitiator    sdk.AccAddress
}

// MakePRNG constructs a PRNG
func MakePRNG(blockHash []byte, blockTimestamp int64, txInitiator sdk.AccAddress) PRNG {
	return PRNG{
		BlockHash:      blockHash,
		BlockTimestamp: blockTimestamp,
		TxInitiator:    txInitiator,
	}
}

// GetRand implements RNG
func (p PRNG) GetRand() sdk.Dec {
	// TODO
	// seed := p.BlockTimestamp + uint64(sdk.Int(p.BlockHash)/Int(p.BlockTimestamp)) + uint64(sdk.Int(sha256.Sum256([]byte(p.TxInitiator))/sdk.Int(p.BlockTimestampt)))
	// return sdk.NewDec(seed - seed/1000*1000)

	return sdk.ZeroDec()
}
