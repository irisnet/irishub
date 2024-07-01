package token

import (
	"encoding/json"
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	tokencli "mods.irisnet.org/modules/token/client/cli"
	tokentypes "mods.irisnet.org/modules/token/types"
	"mods.irisnet.org/simapp"
)

// TxTestSuite is a suite of end-to-end tests for the nft module
type TxTestSuite struct {
	e2e.TestSuite
}

// TestTxCmd tests all tx command in the nft module
func (s *TxTestSuite) TestTxCmd() {
	val := s.Network.Validators[0]
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
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}
	expectedCode := uint32(0)
	txResult := IssueTokenExec(s.T(), s.Network, clientCtx, from.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	tokenSymbol := s.Network.GetAttribute(
		tokentypes.EventTypeIssueToken,
		tokentypes.AttributeKeySymbol,
		txResult.Events,
	)

	//------test GetCmdQueryTokens()-------------
	tokens := QueryTokensExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal(1, len(tokens))

	//------test GetCmdQueryToken()-------------
	token := QueryTokenExec(s.T(), s.Network, clientCtx, tokenSymbol)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(symbol, token.GetSymbol())
	s.Require().Equal(uint64(initialSupply), token.GetInitialSupply())

	//------test GetCmdQueryFee()-------------
	queryFeeResponse := QueryFeeExec(s.T(), s.Network, clientCtx, symbol)
	expectedFeeResp := "{\"exist\":true,\"issue_fee\":{\"denom\":\"stake\",\"amount\":\"13015\"},\"mint_fee\":{\"denom\":\"stake\",\"amount\":\"1301\"}}"
	result, _ := json.Marshal(queryFeeResponse)
	s.Require().Equal(expectedFeeResp, string(result))

	//------test GetCmdQueryParams()-------------
	queryParamsResponse := QueryParamsExec(s.T(), s.Network, clientCtx)
	expectedParams := "{\"token_tax_rate\":\"0.400000000000000000\",\"issue_token_base_fee\":{\"denom\":\"stake\",\"amount\":\"60000\"},\"mint_token_fee_ratio\":\"0.100000000000000000\",\"enable_erc20\":true}"
	result, _ = json.Marshal(queryParamsResponse)
	s.Require().Equal(expectedParams, string(result))

	//------test GetCmdMintToken()-------------
	balance := simapp.QueryBalanceExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		symbol,
	)
	initAmount := balance.Amount.Int64()
	mintAmount := int64(50000000)

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagTo, from.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}
	coinMintedStr := fmt.Sprintf("%d%s", mintAmount, symbol)

	txResult = MintTokenExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		coinMintedStr,
		args...,
	)
	s.Require().Equal(expectedCode, txResult.Code)

	balance = simapp.QueryBalanceExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		symbol,
	)
	exceptedAmount := initAmount + mintAmount
	s.Require().Equal(exceptedAmount, balance.Amount.Int64())

	//------test GetCmdBurnToken()-------------

	burnAmount := int64(2000000)

	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	coinBurntStr := fmt.Sprintf("%d%s", burnAmount, symbol)
	txResult = BurnTokenExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		coinBurntStr,
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	balance = simapp.QueryBalanceExec(
		s.T(),
		s.Network,
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
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = EditTokenExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		symbol,
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	token2 := QueryTokenExec(s.T(), s.Network, clientCtx, tokenSymbol)
	s.Require().Equal(newName, token2.GetName())
	s.Require().Equal(uint64(newMaxSupply), token2.GetMaxSupply())
	s.Require().Equal(newMintable, token2.GetMintable())

	//------test GetCmdTransferTokenOwner()-------------
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagTo, to.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = TransferTokenOwnerExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		symbol,
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	token3 := QueryTokenExec(s.T(), s.Network, clientCtx, tokenSymbol)
	s.Require().Equal(to, token3.GetOwner())
	// ---------------------------------------------------------------------------

	//------test GetCmdSwapToErc20()-------------
	// args = []string{
	// 	fmt.Sprintf("--%s=%s", tokencli.FlagTo, to.String()),

	// 	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	// 	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	// 	fmt.Sprintf(
	// 		"--%s=%s",
	// 		flags.FlagFees,
	// 		sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
	// 	),
	// }

	// txResult = SwapToERC20Exec(
	// 	s.T(),
	// 	s.Network,
	// 	clientCtx,
	// 	from.String(),
	// 	sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(1))).String(),
	// 	args...)

	// TODO assert
	// s.Require().Equal(expectedCode, txResult.Code)
	// ---------------------------------------------------------------------------

	//------test GetCmdSwapFromErc20()-------------
	// args = []string{
	// 	fmt.Sprintf("--%s=%s", tokencli.FlagTo, to.String()),

	// 	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	// 	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	// 	fmt.Sprintf(
	// 		"--%s=%s",
	// 		flags.FlagFees,
	// 		sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
	// 	),
	// }

	// txResult = SwapFromERC20Exec(
	// 	s.T(),
	// 	s.Network,
	// 	clientCtx,
	// 	from.String(),
	// 	sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(1))).String(),
	// 	args...)

	// TODO assert
	// s.Require().Equal(expectedCode, txResult.Code)
	// ---------------------------------------------------------------------------
}
