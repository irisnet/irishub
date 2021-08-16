package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/stretchr/testify/suite"
)

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestInitGenesisAndExportGenesis() {
	expGenesis := types.GenesisState{
		Params: types.Params{
			Fee: sdk.NewDecWithPrec(4, 3),
		},
		StandardDenom: denomStandard,
		Pool: []types.Pool{{
			Id:                types.GetPoolId(denomETH),
			StandardDenom:     denomStandard,
			CounterpartyDenom: denomETH,
			EscrowAddress:     types.GetReservePoolAddr("lpt-1").String(),
			LptDenom:          "lpt-1",
		}},
		Sequence: 2,
	}
	suite.app.CoinswapKeeper.InitGenesis(suite.ctx, expGenesis)
	actGenesis := suite.app.CoinswapKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(expGenesis, actGenesis)
}
