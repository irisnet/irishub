package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/spf13/viper"
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

func TestRestServer(t *testing.T) {
	cdc := app.MakeCodec()
	comm := ServeCommand(cdc)
	viper.Set("chain-id", "fuxi-1001")
	viper.Set("node", "tcp://localhost:26657")
	viper.Set("laddr", "tcp://localhost:1317")

	comm.ExecuteC()
}

func TestMetricsCmd(t *testing.T) {
	cdc := app.MakeCodec()
	comm := prometheus.MonitorCommand(cdc)
	viper.Set("node", "tcp://0.0.0.0:26657")
	comm.ExecuteC()
}

func TestCoins(t *testing.T) {
	someCoins := sdk.Coins{sdk.NewCoin("atom", 123)}
	bz,_ := json.Marshal(someCoins)
	fmt.Println(string(bz))
}

func TestSign(t *testing.T) {
	cdc := app.MakeCodec()

	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	addr2, _ := sdk.AccAddressFromHex("5CBFAEB61B3EA7A17F1D7920EA8E45CEA2C8CC1F")
	coins := sdk.Coins{sdk.NewCoin("iris", 10)}
	var msg = bank.MsgSend{
		Inputs:  []bank.Input{bank.NewInput(addr1, coins)},
		Outputs: []bank.Output{bank.NewOutput(addr2, coins)},
	}

	chainID := "fuxi-10001"
	accnum := int64(0)
	sequence := int64(1)
	fee := auth.NewStdFee(10000, sdk.NewCoin("iris", 0))
	memo := "test"
	msgs := []sdk.Msg{msg}

	signMsg := auth.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: accnum,
		Sequence:      sequence,
		Msgs:          msgs,
		Memo:          memo,
		Fee:           fee, // TODO run simulate to estimate gas?
	}

	signByte1 := signMsg.Bytes()



	db := dbm.NewMemDB()
	cstore := keys.New(db)

	info, err := cstore.CreateKey("iris", "prevent body bachelor lawsuit angry yard squeeze forest young category bargain crash sausage consider amateur next senior burger average state record edit topic abstract", "1234567890")
	if  err != nil {
		panic(err)
	}

	sig, _, _ := cstore.Sign("iris", "1234567890", signByte1)

	sigs := []auth.StdSignature{{
		PubKey:        info.GetPubKey(),
		Signature:     sig,
		AccountNumber: accnum,
		Sequence:      sequence,
	}}

	var stdTx = auth.StdTx{
		Msgs:       signMsg.Msgs,
		Fee:        signMsg.Fee,
		Signatures: sigs,
		Memo:       memo,
	}

	txBytes,_ := cdc.MarshalBinary(stdTx)
	fmt.Println(txBytes)
	//fmt.Println(info.GetPubKey().Address())

	//fmt.Println(info.GetPubKey().Bytes())
	pass := info.GetPubKey().VerifyBytes(signByte1,sig)
	require.Equal(t, true, pass)
	//fmt.Println(sig.Bytes())
}
