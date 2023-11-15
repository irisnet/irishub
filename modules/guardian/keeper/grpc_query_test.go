package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/irisnet/irishub/v2/modules/guardian/types"
)

func (suite *KeeperTestSuite) TestGRPCQuerySupers() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	guardian := types.NewSuper("test", types.Ordinary, addr, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.GuardianKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Supers(gocontext.Background(), &types.QuerySupersRequest{})
	suite.Require().NoError(err)

	app.GuardianKeeper.AddSuper(ctx, guardian)

	supersResp, err := queryClient.Supers(gocontext.Background(), &types.QuerySupersRequest{})
	suite.Require().NoError(err)
	suite.Len(supersResp.Supers, 1)
	suite.Equal(guardian, supersResp.Supers[0])
}
