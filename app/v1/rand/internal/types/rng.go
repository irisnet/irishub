package types

import (
	"math/big"

	sdk "github.com/irisnet/irishub/types"
)

const RandPrec = 20 // the precision for generated random numbers

// RNG is a random number generator
type RNG interface {
	GetRand() sdk.Rat // interface which returns a random number between (0,1)
}

// PRNG represents a pseudo-random number implementation based on future block for RNG
type PRNG struct {
	BlockHash   []byte         // hash of the future block
	TxInitiator sdk.AccAddress // address initiating the request tx
}

// MakePRNG constructs a PRNG
func MakePRNG(blockHash []byte, txInitiator sdk.AccAddress) PRNG {
	return PRNG{
		BlockHash:   blockHash,
		TxInitiator: txInitiator,
	}
}

// GetRand implements RNG
func (p PRNG) GetRand() sdk.Rat {
	seedBH := new(big.Int).SetBytes(sdk.SHA256(p.BlockHash))
	seedTI := new(big.Int).SetBytes(sdk.SHA256(p.TxInitiator))

	seed := new(big.Int).SetBytes(sdk.SHA256(new(big.Int).Add(seedBH, seedTI).Bytes()))
	precision := new(big.Int).Exp(big.NewInt(10), big.NewInt(RandPrec), nil)

	// Generate a random number between [0,1) with `RandPrec` precision from seed
	rand := sdk.NewRatFromBigInt(new(big.Int).Mod(seed, precision), precision)

	return rand
}
