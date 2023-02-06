package cli_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokencli "github.com/irisnet/irismod/modules/token/client/cli"
	tokentestutil "github.com/irisnet/irismod/modules/token/client/testutil"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	network simapp.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = simapp.SetupNetwork(s.T())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestToken() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// ---------------------------------------------------------------------------

	from := val.Address
	symbol := "kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := 0
	initialSupply := int64(100000000)
	maxSupply := int64(200000000)
	mintable := true

	//------test GetCmdIssueToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMinUnit, minUnit),
		fmt.Sprintf("--%s=%d", tokencli.FlagScale, scale),
		fmt.Sprintf("--%s=%d", tokencli.FlagInitialSupply, initialSupply),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}
	expectedCode := uint32(0)
	txResult := tokentestutil.IssueTokenExec(s.T(), s.network, clientCtx, from.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	tokenSymbol := s.network.GetAttribute(tokentypes.EventTypeIssueToken, tokentypes.AttributeKeySymbol, txResult.Events)

	//------test GetCmdQueryTokens()-------------
	tokens := tokentestutil.QueryTokensExec(s.T(), s.network, clientCtx, from.String())
	s.Require().Equal(1, len(tokens))

	//------test GetCmdQueryToken()-------------
	token := tokentestutil.QueryTokenExec(s.T(), s.network, clientCtx, tokenSymbol)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(symbol, token.GetSymbol())
	s.Require().Equal(uint64(initialSupply), token.GetInitialSupply())

	//------test GetCmdQueryFee()-------------
	queryFeeResponse := tokentestutil.QueryFeeExec(s.T(), s.network, clientCtx, symbol)
	expectedFeeResp := "{\"exist\":true,\"issue_fee\":{\"denom\":\"stake\",\"amount\":\"13015\"},\"mint_fee\":{\"denom\":\"stake\",\"amount\":\"1301\"}}"
	result, _ := json.Marshal(queryFeeResponse)
	s.Require().Equal(expectedFeeResp, string(result))

	//------test GetCmdQueryParams()-------------
	queryParamsResponse := tokentestutil.QueryParamsExec(s.T(), s.network, clientCtx)
	expectedParams := "{\"token_tax_rate\":\"0.400000000000000000\",\"issue_token_base_fee\":{\"denom\":\"stake\",\"amount\":\"60000\"},\"mint_token_fee_ratio\":\"0.100000000000000000\"}"
	result, _ = json.Marshal(queryParamsResponse)
	s.Require().Equal(expectedParams, string(result))

	//------test GetCmdMintToken()-------------
	balance := simapp.QueryBalanceExec(
		s.T(),
		s.network,
		clientCtx,
		from.String(),
		symbol,
	)
	initAmount := balance.Amount.Int64()
	mintAmount := int64(50000000)

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagTo, from.String()),
		fmt.Sprintf("--%s=%d", tokencli.FlagAmount, mintAmount),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = tokentestutil.MintTokenExec(s.T(), s.network, clientCtx, from.String(), symbol, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	balance = simapp.QueryBalanceExec(
		s.T(),
		s.network,
		clientCtx,
		from.String(),
		symbol,
	)
	exceptedAmount := initAmount + mintAmount
	s.Require().Equal(exceptedAmount, balance.Amount.Int64())

	//------test GetCmdBurnToken()-------------

	burnAmount := int64(2000000)

	args = []string{
		fmt.Sprintf("--%s=%d", tokencli.FlagAmount, burnAmount),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = tokentestutil.BurnTokenExec(s.T(), s.network, clientCtx, from.String(), symbol, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	balance = simapp.QueryBalanceExec(
		s.T(),
		s.network,
		clientCtx,
		from.String(),
		symbol,
	)
	exceptedAmount = exceptedAmount - burnAmount
	s.Require().Equal(exceptedAmount, balance.Amount.Int64())

	//------test GetCmdEditToken()-------------
	newName := "Wd Token"
	newMaxSupply := 200000000
	newMintable := false

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagName, newName),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, newMaxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, newMintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = tokentestutil.EditTokenExec(s.T(), s.network, clientCtx, from.String(), symbol, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	token2 := tokentestutil.QueryTokenExec(s.T(), s.network, clientCtx, tokenSymbol)
	s.Require().Equal(newName, token2.GetName())
	s.Require().Equal(uint64(newMaxSupply), token2.GetMaxSupply())
	s.Require().Equal(newMintable, token2.GetMintable())

	//------test GetCmdTransferTokenOwner()-------------
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagTo, to.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = tokentestutil.TransferTokenOwnerExec(s.T(), s.network, clientCtx, from.String(), symbol, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	token3 := tokentestutil.QueryTokenExec(s.T(), s.network, clientCtx, tokenSymbol)
	s.Require().Equal(to, token3.GetOwner())
	// ---------------------------------------------------------------------------
}
