package mt

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	mtcli "mods.irisnet.org/modules/mt/client/cli"
	mttypes "mods.irisnet.org/modules/mt/types"
)

// TxTestSuite is a suite of end-to-end tests for the mt module
type TxTestSuite struct {
	e2e.TestSuite
}

// TestMT tests all tx command in the mt module
func (s *TxTestSuite) TestMT() {
	val := s.Validators[0]
	val2 := s.Validators[1]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------
	denomName := "name"
	data := "data"
	from := val.Address
	mintAmt := "10"
	transferAmt := "5"
	burnAmt := "5"

	expectedCode := uint32(0)

	//------test GetCmdIssueDenom()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagName, denomName),
		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult := IssueDenomExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		args...,
	)
	denomID := s.GetAttribute(
		mttypes.EventTypeIssueDenom,
		mttypes.AttributeKeyDenomID,
		txResult.Events,
	)

	//------test GetCmdQueryDenom()-------------
	queryDenomRespType := QueryDenomExec(s.T(), s.Network, clientCtx, denomID)
	s.Require().Equal(denomName, queryDenomRespType.Name)
	s.Require().Equal([]byte(data), queryDenomRespType.Data)

	//------test GetCmdQueryDenoms()-------------
	queryDenomsRespType := QueryDenomsExec(s.T(), s.Network, clientCtx)
	s.Require().Equal(1, len(queryDenomsRespType.Denoms))
	s.Require().Equal(denomID, queryDenomsRespType.Denoms[0].Id)

	//------test GetCmdMintMT()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", mtcli.FlagAmount, mintAmt),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(100))).String(),
		),
	}

	txResult = MintMTExec(s.T(),
		s.Network,
		clientCtx, from.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	mtID := s.GetAttribute(
		mttypes.EventTypeMintMT,
		mttypes.AttributeKeyMTID,
		txResult.Events,
	)
	//------test GetCmdQueryMT()-------------
	queryMTResponse := QueryMTExec(s.T(), s.Network, clientCtx, denomID, mtID)
	s.Require().Equal(mtID, queryMTResponse.Id)

	//-------test GetCmdQueryBalances()----------
	queryBalancesResponse := QueryBlancesExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		denomID,
	)
	s.Require().Equal(1, len(queryBalancesResponse.Balance))
	s.Require().Equal(uint64(10), queryBalancesResponse.Balance[0].Amount)

	//------test GetCmdEditMT()-------------
	newTokenDate := "newdata"
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagData, newTokenDate),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = EditMTExec(s.T(),
		s.Network,
		clientCtx, from.String(), denomID, mtID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryMTResponse = QueryMTExec(s.T(), s.Network, clientCtx, denomID, mtID)
	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)

	//------test GetCmdTransferMT()-------------
	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = TransferMTExec(s.T(),
		s.Network,
		clientCtx, from.String(), recipient.String(), denomID, mtID, transferAmt, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryMTResponse = QueryMTExec(s.T(), s.Network, clientCtx, denomID, mtID)
	s.Require().Equal(mtID, queryMTResponse.Id)
	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)

	//------test GetCmdBurnMT()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = BurnMTExec(s.T(),
		s.Network,
		clientCtx, from.String(), denomID, mtID, burnAmt, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryMTResponse = QueryMTExec(s.T(), s.Network, clientCtx, denomID, mtID)
	s.Require().Equal(mtID, queryMTResponse.Id)
	s.Require().Equal([]byte(newTokenDate), queryMTResponse.Data)
	s.Require().Equal(uint64(5), queryMTResponse.Supply)

	//------test GetCmdTransferDenom()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = TransferDenomExec(s.T(),
		s.Network,
		clientCtx, from.String(), val2.Address.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryDenomResponse := QueryDenomExec(s.T(), s.Network, clientCtx, denomID)
	s.Require().Equal(val2.Address.String(), queryDenomResponse.Owner)
	s.Require().Equal(denomName, queryDenomResponse.Name)
}
