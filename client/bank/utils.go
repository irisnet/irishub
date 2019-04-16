package bank

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

func ConvertAccountCoin(cliCtx context.CLIContext, acc auth.Account) (BaseAccount, error) {
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
	LoosenToken sdk.Coins `json:"loosen_token"`
	BurnedToken sdk.Coins `json:"burned_token"`
	BondedToken sdk.Dec   `json:"bonded_token"`
}
