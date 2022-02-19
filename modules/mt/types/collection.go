package types

import (
	"github.com/irisnet/irismod/modules/mt/exported"
)

// NewCollection creates a new MT Collection
func NewCollection(denom Denom, mts []exported.MT) (c Collection) {
	c.Denom = denom
	for _, mt := range mts {
		c = c.AddMT(mt.(MT))
	}
	return c
}

// AddMT adds an MT to the collection
func (c Collection) AddMT(mt MT) Collection {
	c.MTs = append(c.MTs, mt)
	return c
}

func (c Collection) Supply() int {
	return len(c.MTs)
}

// NewCollection creates a new MT Collection
func NewCollections(c ...Collection) []Collection {
	return append([]Collection{}, c...)
}
