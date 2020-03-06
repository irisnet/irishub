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

// PRNG represents a pseudo-random number implementation based on block for RNG
type PRNG struct {
	BlockHash      []byte         // hash of some block
	BlockTimestamp int64          // timestamp of the next block
	TxInitiator    sdk.AccAddress // address initiating the request tx
	OracleSeed     []byte         // oracle seed
	Oracle         bool           // oracle method
}

// MakePRNG constructs a PRNG
func MakePRNG(blockHash []byte, blockTimestampt int64, txInitiator sdk.AccAddress, oracleSeed []byte, oracle bool) PRNG {
	return PRNG{
		BlockHash:      blockHash,
		BlockTimestamp: blockTimestampt,
		TxInitiator:    txInitiator,
		OracleSeed:     oracleSeed,
		Oracle:         oracle,
	}
}

// GetRand implements RNG
func (p PRNG) GetRand() sdk.Rat {
	seedBT := big.NewInt(p.BlockTimestamp)
	seedBH := new(big.Int).Div(new(big.Int).SetBytes(sdk.SHA256(p.BlockHash)), seedBT)
	seedTI := new(big.Int).Div(new(big.Int).SetBytes(sdk.SHA256(p.TxInitiator)), seedBT)

	seedSum := new(big.Int).Add(seedBT, seedBH)
	seedSum = new(big.Int).Add(seedSum, seedTI)

	if p.Oracle {
		seedOS := new(big.Int).Div(new(big.Int).SetBytes(sdk.SHA256(p.OracleSeed)), seedBT)
		seedSum = new(big.Int).Add(seedSum, seedOS)
	}

	seed := new(big.Int).SetBytes(sdk.SHA256(seedSum.Bytes()))

	precision := new(big.Int).Exp(big.NewInt(10), big.NewInt(RandPrec), nil)

	// Generate a random number between [0,1) with `RandPrec` precision from seed
	rand := sdk.NewRatFromBigInt(new(big.Int).Mod(seed, precision), precision)

	return rand
}
