package types

import abci "github.com/tendermint/tendermint/abci/types"

// Type for querier functions on keepers to implement to handle custom queries
type Querier = func(ctx Context, path []string, req abci.RequestQuery) (res []byte, err Error)

// defines the params for all list queries:
type PaginationParams struct {
	Page uint64
	Size uint16
}

// creates a new PaginationParams
func NewPaginationParams(page uint64, size uint16) PaginationParams {
	if size > 100 {
		size = 100
	}
	return PaginationParams{
		Page: page,
		Size: size,
	}
}

func GetSkipCount(page uint64, size uint16) uint64 {
	if page < 1 {
		page = 1
	}
	return uint64(int(page-1) * int(size))
}

func MarshalErr(err error) Error {
	return ErrInternal(AppendMsgToErr("could not marshal result to JSON", err.Error()))
}
