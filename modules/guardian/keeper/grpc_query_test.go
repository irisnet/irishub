package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/irisnet/irishub/v3/modules/guardian/types"
)

func (suite *KeeperTestSuite) TestGRPCQuerySupers() {
	_, _, addr := testdata.KeyTestPubAddr()
	guardian := types.NewSuper("test", types.Ordinary, addr, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.ifr)
	types.RegisterQueryServer(queryHelper, suite.keeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Supers(gocontext.Background(), &types.QuerySupersRequest{})
	suite.Require().NoError(err)

	suite.keeper.AddSuper(suite.ctx, guardian)
	supersResp, err := queryClient.Supers(gocontext.Background(), &types.QuerySupersRequest{})
	suite.Require().NoError(err)
	suite.Len(supersResp.Supers, 1)
	suite.Equal(guardian, supersResp.Supers[0])
}
