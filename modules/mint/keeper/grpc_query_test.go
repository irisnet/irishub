package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/irisnet/irishub/v2/modules/mint/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryPoolParameters() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.MintKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(app.MintKeeper.GetParams(ctx), resp.Params)
}
