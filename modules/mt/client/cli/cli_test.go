package cli_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/cometbft/cometbft/crypto"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/simapp"
// 	mtcli "github.com/irisnet/irismod/mt/client/cli"
// 	mttestutil "github.com/irisnet/irismod/mt/client/testutil"
// 	mttypes "github.com/irisnet/irismod/mt/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestMT() {
// 	val := s.network.Validators[0]
// 	val2 := s.network.Validators[1]
// 	clientCtx := val.ClientCtx

// 	// ---------------------------------------------------------------------------
// 	denomName := "name"
// 	data := "data"
// 	from := val.Address
// 	mintAmt := "10"
// 	transferAmt := "5"
// 	burnAmt := "5"

// 	expectedCode := uint32(0)

// 	//------test GetCmdIssueDenom()-------------
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", mtcli.FlagName, denomName),
// 		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult := mttestutil.IssueDenomExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		from.String(),
// 		args...,
// 	)
// 	denomID := s.network.GetAttribute(
// 		mttypes.EventTypeIssueDenom,
// 		mttypes.AttributeKeyDenomID,
// 		txResult.Events,
// 	)

// 	//------test GetCmdQueryDenom()-------------
// 	queryDenomRespType := mttestutil.QueryDenomExec(s.T(), s.network, clientCtx, denomID)
// 	s.Require().Equal(denomName, queryDenomRespType.Name)
// 	s.Require().Equal([]byte(data), queryDenomRespType.Data)

// 	//------test GetCmdQueryDenoms()-------------
// 	queryDenomsRespType := mttestutil.QueryDenomsExec(s.T(), s.network, clientCtx)
// 	s.Require().Equal(1, len(queryDenomsRespType.Denoms))
// 	s.Require().Equal(denomID, queryDenomsRespType.Denoms[0].Id)

// 	//------test GetCmdMintMT()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
// 		fmt.Sprintf("--%s=%s", mtcli.FlagAmount, mintAmt),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(100))).String(),
// 		),
// 	}

// 	txResult = mttestutil.MintMTExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), denomID, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	mtID := s.network.GetAttribute(
// 		mttypes.EventTypeMintMT,
// 		mttypes.AttributeKeyMTID,
// 		txResult.Events,
// 	)
// 	//------test GetCmdQueryMT()-------------
// 	queryMTResponse := mttestutil.QueryMTExec(s.T(), s.network, clientCtx, denomID, mtID)
// 	s.Require().Equal(mtID, queryMTResponse.Id)

// 	//-------test GetCmdQueryBalances()----------
// 	queryBalancesResponse := mttestutil.QueryBlancesExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		from.String(),
// 		denomID,
// 	)
// 	s.Require().Equal(1, len(queryBalancesResponse.Balance))
// 	s.Require().Equal(uint64(10), queryBalancesResponse.Balance[0].Amount)

// 	//------test GetCmdEditMT()-------------
// 	newTokenDate := "newdata"
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", mtcli.FlagData, newTokenDate),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = mttestutil.EditMTExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), denomID, mtID, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	queryMTResponse = mttestutil.QueryMTExec(s.T(), s.network, clientCtx, denomID, mtID)
// 	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)

// 	//------test GetCmdTransferMT()-------------
// 	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = mttestutil.TransferMTExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), recipient.String(), denomID, mtID, transferAmt, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	queryMTResponse = mttestutil.QueryMTExec(s.T(), s.network, clientCtx, denomID, mtID)
// 	s.Require().Equal(mtID, queryMTResponse.Id)
// 	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)

// 	//------test GetCmdBurnMT()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = mttestutil.BurnMTExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), denomID, mtID, burnAmt, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	queryMTResponse = mttestutil.QueryMTExec(s.T(), s.network, clientCtx, denomID, mtID)
// 	s.Require().Equal(mtID, queryMTResponse.Id)
// 	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)
// 	s.Require().Equal(uint64(5), queryMTResponse.Supply)

// 	//------test GetCmdTransferDenom()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = mttestutil.TransferDenomExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), val2.Address.String(), denomID, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	queryDenomResponse := mttestutil.QueryDenomExec(s.T(), s.network, clientCtx, denomID)
// 	s.Require().Equal(val2.Address.String(), queryDenomResponse.Owner)
// 	s.Require().Equal(denomName, queryDenomResponse.Name)
// }
