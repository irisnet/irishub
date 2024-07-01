package types

import (
	"crypto/sha256"
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RandPrec        = 20 // the precision for generated random numbers
	SeedBytesLength = 32 //
)

// RNG is a random number generator
type RNG interface {
	GetRand() big.Rat // interface which returns a random number between [0,1)
}

// PRNG represents a pseudo-random number implementation based on block or oracle for RNG
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
func (p PRNG) GetRand() *big.Rat {
	seedBT := big.NewInt(p.BlockTimestamp)
	seedBH := new(big.Int).Div(new(big.Int).SetBytes(SHA256(p.BlockHash)), seedBT)
	seedTI := new(big.Int).Div(new(big.Int).SetBytes(SHA256(p.TxInitiator)), seedBT)

	seedSum := new(big.Int).Add(seedBT, seedBH)
	seedSum = new(big.Int).Add(seedSum, seedTI)

	if p.Oracle {
		seedOS := new(big.Int).Div(new(big.Int).SetBytes(SHA256(p.OracleSeed)), seedBT)
		seedSum = new(big.Int).Add(seedSum, seedOS)
	}

	seed := new(big.Int).SetBytes(SHA256(seedSum.Bytes()))

	precision := new(big.Int).Exp(big.NewInt(10), big.NewInt(RandPrec), nil)

	// Generate a random number between [0,1) with `RandPrec` precision from seed
	random := new(big.Rat).SetFrac(new(big.Int).Mod(seed, precision), precision)

	return random
}

func (p PRNG) Intn(n int) int {
	rnd := rand.New(rand.NewSource(p.GetRand().Denom().Int64()))
	return rnd.Intn(n)
}

// SHA256 wraps sha256.Sum256 with result converted to slice
func SHA256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
