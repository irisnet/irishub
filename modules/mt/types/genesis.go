package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/irisnet/irismod/mt/exported"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(collections []Collection, owners []Owner) *GenesisState {
	return &GenesisState{
		Collections: collections,
		Owners:      owners,
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

	mtCount1 := len(mtMap1)

	// --------------------------
	var denomMap2 map[string]bool
	denomMap2 = make(map[string]bool)

	var mtMap2 map[string]uint64
	mtMap2 = make(map[string]uint64)

	for _, o := range data.Owners {
		for _, d := range o.Denoms {
			denomMap2[d.DenomId] = true

			if _, ok := denomMap1[d.DenomId]; !ok {
				return errorsmod.Wrapf(errorsmod.ErrPanic, "unknown mt denom, (%s)", d.DenomId)
			}
			for _, b := range d.Balances {
				mtMap2[d.DenomId+b.MtId] = mtMap2[d.DenomId+b.MtId] + b.Amount

			}
		}
	}

	mtCount2 := len(mtMap2)

	if mtCount1 != mtCount2 {
		return errorsmod.Wrapf(errorsmod.ErrPanic, "mt count mismatch, (%d, %d)", mtCount1, mtCount2)
	}

	for id1, supply1 := range mtMap1 {
		supply2 := mtMap2[id1]
		if supply1 != supply2 {
			return errorsmod.Wrapf(errorsmod.ErrPanic, "mt supply mismatch, id: %s (%d, %d)", id1, supply1, mtCount2)
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
