package bank

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/client/context"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
  Address:         %s
  Pubkey:          %s
  Coins:           %s
  Account Number:  %d
  Sequence:        %d`,
		acc.Address, pubkey, strings.Join(acc.Coins, ","), acc.AccountNumber, acc.Sequence,
	)
}

func ConvertToMainUnit(cliCtx context.CLIContext, coins sdk.Coins) (resCoins []string, err error) {
	for _, coin := range coins {
		mainUnit, err := sdk.GetCoinName(coin.String())
		if err != nil {
			return nil, err
		}

		ct, err := cliCtx.GetCoinType(mainUnit)
		if err != nil {
			return nil, err
		}

		destCoinStr, err := ct.Convert(coin.String(), mainUnit)
		if err != nil {
			return nil, err
		}

		resCoins = append(resCoins, destCoinStr)
	}

	return resCoins, nil
}

func ConvertAccountCoin(cliCtx context.CLIContext, acc auth.BaseAccount) (BaseAccount, error) {
	var accCoins []string
	for _, coin := range acc.GetCoins() {
		coinString, err := cliCtx.ConvertCoinToMainUnit(coin.String())
		if err == nil {
			accCoins = append(accCoins, coinString[0])
		} else {
			accCoins = append(accCoins, coin.String())
		}

	}
	return BaseAccount{
		Address:       acc.GetAddress(),
		Coins:         accCoins,
		PubKey:        acc.GetPubKey(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}, nil
}

// BuildBankSendMsg builds the sending coins msg
func BuildBankSendMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}

// BuildBankBurnMsg builds the burning coin msg
func BuildBankBurnMsg(from sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	msg := bank.NewMsgBurn(from, coins)
	return msg
}

type TokenStats struct {
	LooseTokens  []string `json:"loose_tokens"`
	BurnedTokens []string `json:"burned_tokens"`
	BondedTokens []string `json:"bonded_tokens"`
	TotalSupply  []string `json:"total_supply"`
}

// String implements fmt.Stringer
func (ts TokenStats) String() string {
	return fmt.Sprintf(`TokenStats:
  Loose Tokens:   %s
  Bonded Tokens:  %s
  Burned Tokens:  %s
  Total Supply:   %s`,
		strings.Join(ts.LooseTokens, ","), strings.Join(ts.BondedTokens, ","), strings.Join(ts.BurnedTokens, ","), strings.Join(ts.TotalSupply, ","),
	)
}
