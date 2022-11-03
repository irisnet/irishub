package keeper

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"testing"
)

var ShapePageRequest = shapePageRequest

func TestShapePageRequest(t *testing.T) {
	defaultRequest := newDefaultPageRequest()

	res1 := ShapePageRequest(nil)
	require.NotNil(t, res1)
	require.Equal(t, res1.Limit, defaultRequest.Limit)           // limit == paginationDefaultLimit
	require.Equal(t, res1.Offset, defaultRequest.Offset)         // offset == 0
	require.Equal(t, res1.CountTotal, defaultRequest.CountTotal) // count_total == false

	request := &query.PageRequest{
		Key:        nil,
		Offset:     100,
		Limit:      10000,
		CountTotal: true,
		Reverse:    true,
	}

	res2 := ShapePageRequest(request)
	require.NotNil(t, res2)
	require.Equal(t, res2.Limit, defaultRequest.Limit)           // limit == paginationDefaultLimit
	require.Equal(t, res2.Offset, defaultRequest.Offset)         // offset == 0
	require.Equal(t, res2.CountTotal, defaultRequest.CountTotal) // count_total == false
}
