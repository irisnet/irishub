package app

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gogo/protobuf/proto"

	convertertypes "github.com/irisnet/erc721-bridge/x/converter/types"
)

func QueryTokenTrace(
	ctx client.Context,
	classId, tokenId string,
) (proto.Message, error) {
	queryClient := convertertypes.NewQueryClient(ctx)
	req := &convertertypes.QueryTokenTraceRequest{
		ClassId: classId,
		TokenId: tokenId,
	}
	return queryClient.TokenTrace(context.Background(), req)
}
