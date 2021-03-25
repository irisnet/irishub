package htlc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

type GenesisTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	app    *simapp.SimApp
	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *GenesisTestSuite) SetupTest() {
	app := simapp.Setup(false)
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.keeper = &app.HTLCKeeper
	suite.app = app

	_, addrs := GeneratePrivKeyAddressPairs(3)
	suite.addrs = addrs
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestGenesisState() {
	type GenState func() *types.GenesisState

	testCases := []struct {
		name       string
		genState   GenState
		expectPass bool
	}{{
		name: "default",
		genState: func() *types.GenesisState {
			return NewHTLTGenesis(suite.addrs[0])
		},
		expectPass: true,
	}, {
		name: "import atomic htlcs and asset supplies",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			_, addrs := GeneratePrivKeyAddressPairs(2)
			var htlcs []types.HTLC
			var supplies []types.AssetSupply
			for i := 0; i < 2; i++ {
				htlc, supply := loadSwapAndSupply(addrs[i], i)
				htlcs = append(htlcs, htlc)
				supplies = append(supplies, supply)
			}
			gs.Htlcs = htlcs
			gs.Supplies = supplies
			return gs
		},
		expectPass: true,
	}, {
		name: "0 deputy fees",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Params.AssetParams[0].FixedFee = sdk.ZeroInt()
			return gs
		},
		expectPass: true,
	}, {
		name: "incoming supply doesn't match amount in incoming atomic swaps",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			_, addrs := GeneratePrivKeyAddressPairs(1)
			swap, _ := loadSwapAndSupply(addrs[0], 2)
			gs.Htlcs = []types.HTLC{swap}
			return gs
		},
		expectPass: false,
	}, {
		name: "current supply above limit",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			assetParam, _ := suite.keeper.GetAsset(suite.ctx, "htltbnb")
			gs.Supplies = []types.AssetSupply{
				{
					IncomingSupply: c("htltbnb", 0),
					OutgoingSupply: c("htltbnb", 0),
					CurrentSupply:  c("htltbnb", assetParam.SupplyLimit.Limit.Add(i(1)).Int64()),
				},
			}
			return gs
		},
		expectPass: false,
	}, {
		name: "incoming supply above limit",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])

			assetParam, _ := suite.keeper.GetAsset(suite.ctx, "htltbnb")
			overLimitAmount := assetParam.SupplyLimit.Limit.Add(i(1))

			_, addrs := GeneratePrivKeyAddressPairs(2)
			timestamp := ts(0)
			randomSecret, _ := GenerateRandomSecret()
			randomHashLock := types.GetHashLock(randomSecret, timestamp)
			amount := cs(c("htltbnb", overLimitAmount.Int64()))
			id := types.GetID(suite.addrs[0], addrs[1], amount, randomHashLock)

			htlc := types.NewHTLC(
				id,
				suite.addrs[0],
				addrs[1],
				ReceiverOnOtherChain,
				SenderOnOtherChain,
				amount,
				randomHashLock,
				[]byte{},
				timestamp,
				MaxTimeLock,
				types.Open,
				0,
				true,
				types.Incoming,
			)
			gs.Htlcs = []types.HTLC{htlc}

			gs.Supplies = []types.AssetSupply{
				{
					IncomingSupply: c("htltbnb", assetParam.SupplyLimit.Limit.Add(i(1)).Int64()),
					OutgoingSupply: c("htltbnb", 0),
					CurrentSupply:  c("htltbnb", 0),
				},
			}
			return gs
		},
		expectPass: false,
	}, {
		name: "incoming supply + current supply above limit",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])

			assetParam, _ := suite.keeper.GetAsset(suite.ctx, "htltbnb")
			halfLimit := assetParam.SupplyLimit.Limit.Int64() / 2
			overHalfLimit := halfLimit + 1

			_, addrs := GeneratePrivKeyAddressPairs(2)
			timestamp := ts(0)
			randomSecret, _ := GenerateRandomSecret()
			randomHashLock := types.GetHashLock(randomSecret, timestamp)
			amount := cs(c("htltbnb", halfLimit))
			id := types.GetID(suite.addrs[0], addrs[1], amount, randomHashLock)

			htlc := types.NewHTLC(
				id,
				suite.addrs[0],
				addrs[1],
				ReceiverOnOtherChain,
				SenderOnOtherChain,
				amount,
				randomHashLock,
				[]byte{},
				timestamp,
				uint64(360),
				types.Open,
				0,
				true,
				types.Incoming,
			)
			gs.Htlcs = []types.HTLC{htlc}

			gs.Supplies = []types.AssetSupply{
				{
					IncomingSupply: c("htltbnb", halfLimit),
					OutgoingSupply: c("htltbnb", 0),
					CurrentSupply:  c("htltbnb", overHalfLimit),
				},
			}
			return gs
		},
		expectPass: false,
	}, {
		name: "outgoing supply above limit",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])

			assetParam, _ := suite.keeper.GetAsset(suite.ctx, "htltbnb")
			overLimitAmount := assetParam.SupplyLimit.Limit.Add(i(1))

			_, addrs := GeneratePrivKeyAddressPairs(2)
			timestamp := ts(0)
			randomSecret, _ := GenerateRandomSecret()
			randomHashLock := types.GetHashLock(randomSecret, timestamp)
			amount := cs(c("htltbnb", overLimitAmount.Int64()))
			id := types.GetID(suite.addrs[0], addrs[1], amount, randomHashLock)

			htlc := types.NewHTLC(
				id,
				addrs[1],
				suite.addrs[0],
				ReceiverOnOtherChain,
				SenderOnOtherChain,
				amount,
				randomHashLock,
				[]byte{},
				timestamp,
				MinTimeLock,
				types.Open,
				0,
				true,
				types.Outgoing,
			)
			gs.Htlcs = []types.HTLC{htlc}

			gs.Supplies = []types.AssetSupply{
				{
					IncomingSupply: c("htltbnb", 0),
					OutgoingSupply: c("htltbnb", 0),
					CurrentSupply:  c("htltbnb", assetParam.SupplyLimit.Limit.Add(i(1)).Int64()),
				},
			}
			return gs
		},
		expectPass: false,
	}, {
		name: "asset supply denom is not a supported asset",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Supplies = []types.AssetSupply{
				{
					IncomingSupply: c("fake", 0),
					OutgoingSupply: c("fake", 0),
					CurrentSupply:  c("fake", 0),
				},
			}
			return gs
		},
		expectPass: false,
	}, {
		name: "atomic swap asset type is unsupported",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			_, addrs := GeneratePrivKeyAddressPairs(2)
			timestamp := ts(0)
			randomSecret, _ := GenerateRandomSecret()
			randomHashLock := types.GetHashLock(randomSecret, timestamp)
			amount := cs(c("fake", 500000))
			id := types.GetID(suite.addrs[0], addrs[1], amount, randomHashLock)

			htlc := types.NewHTLC(
				id,
				suite.addrs[0],
				addrs[1],
				ReceiverOnOtherChain,
				SenderOnOtherChain,
				amount,
				randomHashLock,
				[]byte{},
				timestamp,
				uint64(360),
				types.Open,
				0,
				true,
				types.Incoming,
			)
			gs.Htlcs = []types.HTLC{htlc}
			return gs
		},
		expectPass: false,
	}, {
		name: "atomic swap status is invalid",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			_, addrs := GeneratePrivKeyAddressPairs(2)
			timestamp := ts(0)
			randomSecret, _ := GenerateRandomSecret()
			randomHashLock := types.GetHashLock(randomSecret, timestamp)
			amount := cs(c("htltbnb", 5000))
			id := types.GetID(suite.addrs[0], addrs[1], amount, randomHashLock)

			htlc := types.NewHTLC(
				id,
				suite.addrs[0],
				addrs[1],
				ReceiverOnOtherChain,
				SenderOnOtherChain,
				amount,
				randomHashLock,
				[]byte{},
				timestamp,
				uint64(360),
				3,
				0,
				true,
				types.Incoming,
			)
			gs.Htlcs = []types.HTLC{htlc}
			return gs
		},
		expectPass: false,
	}, {
		name: "minimum block lock cannot be > maximum block lock",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Params.AssetParams[0].MinBlockLock = 201
			gs.Params.AssetParams[0].MaxBlockLock = 200
			return gs
		},
		expectPass: false,
	}, {
		name: "empty supported asset denom",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Params.AssetParams[0].Denom = ""
			return gs
		},
		expectPass: false,
	}, {
		name: "negative supported asset limit",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Params.AssetParams[0].SupplyLimit.Limit = i(-100)
			return gs
		},
		expectPass: false,
	}, {
		name: "duplicate supported asset denom",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			gs.Params.AssetParams[1].Denom = "htltbnb"
			return gs
		},
		expectPass: false,
	}}

	for _, tc := range testCases {
		suite.Run(
			tc.name,
			func() {
				if tc.expectPass {
					suite.NotPanics(
						func() {
							simapp.SetupWithGenesisHTLC(tc.genState())
						},
						tc.name,
					)
				} else {
					suite.Panics(
						func() {
							simapp.SetupWithGenesisHTLC(tc.genState())
						},
						tc.name,
					)
				}
			},
		)
	}
}
