package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/mt/exported"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(collections []Collection) *GenesisState {
	return &GenesisState{
		Collections: collections,
	}
}

// ValidateGenesis performs basic validation of mts genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {

	var denomMap1 map[string]int
	denomMap1 = make(map[string]int)

	var mtMap1 map[string]uint64
	mtMap1 = make(map[string]uint64)

	for _, c := range data.Collections {
		denomMap1[c.Denom.Id] = len(c.Mts)

		for _, m := range c.Mts {
			mtMap1[c.Denom.Id+m.Id] = m.Supply
		}
	}

	denomCount1 := len(data.Collections)
	mtCount1 := len(mtMap1)

	// --------------------------
	var denomMap2 map[string]bool
	denomMap2 = make(map[string]bool)

	var mtMap2 map[string]uint64
	mtMap2 = make(map[string]uint64)

	for _, o := range data.Owners {
		for _, d := range o.Denoms {
			denomMap2[d.DenomId] = true
			for _, b := range d.Balances {
				mtMap2[d.DenomId+b.MtId] = mtMap2[d.DenomId+b.MtId] + b.Amount

			}
		}
	}

	denomCount2 := len(denomMap2)
	mtCount2 := len(mtMap2)

	if denomCount1 != denomCount2 {
		return sdkerrors.Wrapf(sdkerrors.ErrPanic, "mt denom count mismatch, (%d, %d)", denomCount1, denomCount2)
	}

	if mtCount1 != mtCount2 {
		return sdkerrors.Wrapf(sdkerrors.ErrPanic, "mt count mismatch, (%d, %d)", mtCount1, mtCount2)
	}

	for id1, supply1 := range mtMap1 {
		supply2 := mtMap2[id1]
		if supply1 != supply2 {
			return sdkerrors.Wrapf(sdkerrors.ErrPanic, "mt supply mismatch, id: %s (%d, %d)", id1, supply1, mtCount2)
		}
	}

	return nil
}

// ---------------------------------------------
// NewCollection creates a new MT Collection
func NewCollection(denom Denom, mts []exported.MT) (c Collection) {
	c.Denom = &denom
	for _, mt := range mts {
		c = c.AddMT(mt.(MT))
	}
	return c
}

// AddMT adds an MT to the collection
func (c Collection) AddMT(mt MT) Collection {
	c.Mts = append(c.Mts, mt)
	return c
}

// Supply queries supply of a collection
func (c Collection) Supply() int {
	return len(c.Mts)
}

// NewCollections creates collections
func NewCollections(c ...Collection) []Collection {
	return append([]Collection{}, c...)
}

// ---------------------------------------------
// NewOwner creates a new owner balance
func NewOwner(address string, denoms []DenomBalance) (o Owner) {
	o.Address = address
	o.Denoms = denoms
	return o
}

// NewDenomBalance creates a new denom balance
func NewDenomBalance(denomID string, balances []Balance) (d DenomBalance) {
	d.DenomId = denomID
	d.Balances = balances
	return d
}

// NewBalance creates a new mt balance
func NewBalance(mtID string, amount uint64) (b Balance) {
	b.MtId = mtID
	b.Amount = amount
	return b
}
