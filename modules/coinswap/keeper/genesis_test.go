package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/coinswap/types"
)

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestInitGenesisAndExportGenesis() {
	expGenesis := types.GenesisState{
		Params:        types.DefaultParams(),
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
	suite.keeper.InitGenesis(suite.ctx, expGenesis)
	actGenesis := suite.keeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(expGenesis, actGenesis)
}
