package rand

import (
	"fmt"
	"math/big"

	sdk "github.com/irisnet/irishub/types"
)

const RandPrec = 10 // the precision for generated random numbers. Must be LTE the precision of sdk.Dec

// RNG is a random number generator
type RNG interface {
	GetRand() sdk.Dec // interface which returns a random number between (0,1)
}

// PRNG represents a pseudo-random number implementation based on future block for RNG
type PRNG struct {
	BlockHash      []byte         // hash of the future block
	BlockTimestamp int64          // timestamp of the future block
	TxInitiator    sdk.AccAddress // address initiating the request tx
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
	seedBT := sdk.NewInt(p.BlockTimestamp)
	seedBH := sdk.NewIntFromBigInt(new(big.Int).SetBytes(sdk.SHA256(p.BlockHash))).Div(seedBT)
	seedTI := sdk.NewIntFromBigInt(new(big.Int).SetBytes(sdk.SHA256(p.TxInitiator))).Div(seedBT)

	seed := sdk.NewIntFromBigInt(new(big.Int).SetBytes(sdk.SHA256(seedBT.Add(seedBH).Add(seedTI).BigInt().Bytes())))
	precision := sdk.NewIntWithDecimal(1, RandPrec)

	// Generate a random number with `RandPrec` precision between (0,1).
	// Err will not occur.
	rand, _ := sdk.NewDecFromStr(fmt.Sprintf("0.%s", seed.Sub(seed.Div(precision).Mul(precision))))

	return rand
}
