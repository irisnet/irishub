package types

import (
	"encoding/hex"
	"fmt"
)

type Seed struct {
	Seed []byte `json:"seed"` // oracle seed
}

// NewRand constructs a Rand
func NewSeed(seed []byte) Seed {
	return Seed{Seed: seed}
}

// String implements fmt.Stringer
func (s Seed) String() string {
	return fmt.Sprintf("Seed: %s", hex.EncodeToString(s.Seed))
}
