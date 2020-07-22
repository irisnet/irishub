package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/irisnet/irishub/modules/guardian/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryProfilers() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	guardian := types.NewGuardian("test", types.Ordinary, addr, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.GuardianKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Profilers(gocontext.Background(), &types.QueryProfilersRequest{})
	suite.Require().NoError(err)

	app.GuardianKeeper.AddProfiler(ctx, guardian)

	profilersResp, err := queryClient.Profilers(gocontext.Background(), &types.QueryProfilersRequest{})
	suite.Require().NoError(err)
	suite.Len(profilersResp.Profilers, 1)
	suite.Equal(guardian, profilersResp.Profilers[0])
}

func (suite *KeeperTestSuite) TestGRPCQueryTrustees() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	guardian := types.NewGuardian("test", types.Ordinary, addr, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.GuardianKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Trustees(gocontext.Background(), &types.QueryTrusteesRequest{})
	suite.Require().NoError(err)

	app.GuardianKeeper.AddTrustee(ctx, guardian)

	trusteesResp, err := queryClient.Trustees(gocontext.Background(), &types.QueryTrusteesRequest{})
	suite.Require().NoError(err)
	suite.Len(trusteesResp.Trustees, 1)
	suite.Equal(guardian, trusteesResp.Trustees[0])
}
