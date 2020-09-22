package types

import (
	"bytes"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewIDCollection creates a new IDCollection instance
func NewIDCollection(denom string, ids []string) IDCollection {
	return IDCollection{
		Denom: strings.TrimSpace(denom),
		Ids:   ids,
	}
}

// Supply return the amount of the denom
func (idc IDCollection) Supply() int {
	return len(idc.Ids)
}

// AddID adds an ID to the idCollection
func (idc IDCollection) AddID(id string) IDCollection {
	idc.Ids = append(idc.Ids, id)
	return idc
}

// ----------------------------------------------------------------------------
// IDCollections is an array of ID Collections
type IDCollections []IDCollection

// Add adds an ID to the idCollection
func (idcs IDCollections) Add(denom, id string) IDCollections {
	for i, idc := range idcs {
		if idc.Denom == denom {
			idcs[i] = idc.AddID(id)
			return idcs
		}
	}
	return append(idcs, IDCollection{
		Denom: denom,
		Ids:   []string{id},
	})
}

// String follows stringer interface
func (idcs IDCollections) String() string {
	if len(idcs) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, idCollection := range idcs {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(idCollection.String())
	}
	return buf.String()
}

// Owner of non fungible tokens
//type Owner struct {
//	Address       sdk.AccAddress `json:"address" yaml:"address"`
//	IDCollections IDCollections  `json:"id_collections" yaml:"id_collections"`
//}

// NewOwner creates a new Owner
func NewOwner(owner sdk.AccAddress, idCollections ...IDCollection) Owner {
	return Owner{
		Address:       owner,
		IDCollections: idCollections,
	}
}

type Owners []Owner

// NewOwner creates a new Owner
func NewOwners(owner ...Owner) Owners {
	return append([]Owner{}, owner...)
}

// String follows stringer interface
func (owners Owners) String() string {
	var buf bytes.Buffer
	for _, owner := range owners {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(owner.String())
	}
	return buf.String()
}
