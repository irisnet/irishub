package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/mt/exported"
)

var _ exported.MT = MT{}

// NewMT creates a new MT instance
func NewMT(id string, supply uint64, owner sdk.AccAddress, data []byte) MT {
	return MT{
		Id:     id,
		Supply: supply,
		Owner:  owner.String(),
		Data:   data,
	}
}

// GetID return the id of MT
func (mt MT) GetID() string {
	return mt.Id
}

// GetID return the supply of MT
func (mt MT) GetSupply() uint64 {
	return mt.Supply
}

// GetOwner return the owner of MT
func (mt MT) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(mt.Owner)
	return owner
}

// GetData return the Data of MT
func (mt MT) GetData() []byte {
	return mt.Data
}

// ----------------------------------------------------------------------------
// MT

// MTs define a list of MT
type MTs []exported.MT

// NewMTs creates a new set of MTs
func NewMTs(nfts ...exported.MT) MTs {
	if len(nfts) == 0 {
		return MTs{}
	}
	return MTs(nfts)
}
