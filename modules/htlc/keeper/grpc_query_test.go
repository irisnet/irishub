package keeper_test

// import (
// 	gocontext "context"

// 	"github.com/cosmos/cosmos-sdk/baseapp"

// 	"github.com/irisnet/irismod/modules/htlc/types"
// )

// func (suite *KeeperTestSuite) TestQueryHTLC() {
// 	err := suite.keeper.CreateHTLC(suite.ctx, sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
// 	suite.NoError(err)

// 	app, ctx := suite.app, suite.ctx

// 	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
// 	types.RegisterQueryServer(queryHelper, app.HTLCKeeper)
// 	queryClient := types.NewQueryClient(queryHelper)

// 	response, err := queryClient.HTLC(gocontext.Background(), &types.QueryHTLCRequest{Id: id.String()})

// 	suite.NoError(err)
// 	suite.Equal(amount, response.Htlc.Amount)
// }

// func (suite *KeeperTestSuite) TestQueryAssetSupply() {
// 	// TODO
// }

// func (suite *KeeperTestSuite) TestQueryAssetSupplies() {
// 	// TODO
// }

// func (suite *KeeperTestSuite) TestQueryParams() {
// 	// TODO
// }
