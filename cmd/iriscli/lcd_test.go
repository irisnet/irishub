package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/spf13/viper"
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/wire"
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
	bz := []byte{123,34,97,99,99,111,117,110,116,95,110,117,109,98,101,114,34,58,34,48,34,44,34,99,104,97,105,110,95,105,100,34,58,34,102,117,120,105,45,49,48,48,49,34,44,34,102,101,101,34,58,123,34,97,109,111,117,110,116,34,58,91,123,34,97,109,111,117,110,116,34,58,34,48,34,44,34,100,101,110,111,109,34,58,34,105,114,105,115,34,125,93,44,34,103,97,115,34,58,34,49,48,48,48,48,34,125,44,34,109,101,109,111,34,58,34,116,101,115,116,34,44,34,109,115,103,115,34,58,91,123,34,116,121,112,101,34,58,34,99,111,115,109,111,115,45,115,100,107,47,77,115,103,68,101,108,101,103,97,116,101,34,44,34,118,97,108,117,101,34,58,123,34,100,101,108,101,103,97,116,105,111,110,34,58,123,34,97,109,111,117,110,116,34,58,34,49,48,34,44,34,100,101,110,111,109,34,58,34,105,114,105,115,34,125,44,34,100,101,108,101,103,97,116,111,114,95,97,100,100,114,34,58,34,99,111,115,109,111,115,97,99,99,97,100,100,114,49,109,103,116,121,122,104,118,102,106,52,50,52,113,57,114,51,53,100,109,114,55,48,113,116,116,54,51,99,112,99,53,56,100,99,116,55,106,113,34,44,34,118,97,108,105,100,97,116,111,114,95,97,100,100,114,34,58,34,99,111,115,109,111,115,97,99,99,97,100,100,114,49,109,103,116,121,122,104,118,102,106,52,50,52,113,57,114,51,53,100,109,114,55,48,113,116,116,54,51,99,112,99,53,56,100,99,116,55,106,113,34,125,125,93,44,34,115,101,113,117,101,110,99,101,34,58,34,57,34,125}
	fmt.Println(string(bz))
}


func TestMsgDelegate(t *testing.T) {
	//cdc := app.MakeCodec()
	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	coins := sdk.NewCoin("iris", 10)
	var delegate = stake.NewMsgDelegate(addr1,addr1,coins)
	var sig2 = delegate.GetSignBytes()
	fmt.Println(string(sig2))
	fmt.Println(sig2)
}


func makeDelegateTx(cdc *wire.Codec) sdk.Msg{
	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	coins := sdk.NewCoin("iris", 10)
	var delegate = stake.NewMsgDelegate(addr1,addr1,coins)
	return delegate
}

func makeSendTx(cdc *wire.Codec) sdk.Msg{
	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	addr2, _ := sdk.AccAddressFromHex("5CBFAEB61B3EA7A17F1D7920EA8E45CEA2C8CC1F")
	coins := sdk.Coins{sdk.NewCoin("iris", 10)}
	var msg = bank.MsgSend{
		Inputs:  []bank.Input{bank.NewInput(addr1, coins)},
		Outputs: []bank.Output{bank.NewOutput(addr2, coins)},
	}
	return msg
}

func TestSign(t *testing.T) {
	cdc := app.MakeCodec()
	chainID := "fuxi-1001"
	accnum := int64(0)
	sequence := int64(9)
	fee := auth.NewStdFee(10000, sdk.NewCoin("iris", 0))
	memo := "test"

	msgs := []sdk.Msg{
		//makeSendTx(cdc),
		makeDelegateTx(cdc),
	}

	signMsg := auth.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: accnum,
		Sequence:      sequence,
		Msgs:          msgs,
		Memo:          memo,
		Fee:           fee, // TODO run simulate to estimate gas?
	}

	signByte1 := signMsg.Bytes()

	fmt.Println(signByte1)
	fmt.Println(string(signByte1))

	db := dbm.NewMemDB()
	cstore := keys.New(db)

	info, err := cstore.CreateKey("iris", "prevent body bachelor lawsuit angry yard squeeze forest young category bargain crash sausage consider amateur next senior burger average state record edit topic abstract", "1234567890")
	if  err != nil {
		panic(err)
	}

	sig, _, _ := cstore.Sign("iris", "1234567890", signByte1)

	//sigs := []auth.StdSignature{{
	//	PubKey:        info.GetPubKey(),
	//	Signature:     sig,
	//	AccountNumber: accnum,
	//	Sequence:      sequence,
	//}}

	//var stdTx = auth.StdTx{
	//	Msgs:       signMsg.Msgs,
	//	Fee:        signMsg.Fee,
	//	Signatures: sigs,
	//	Memo:       memo,
	//}

	//txBytes,_ := cdc.MarshalBinary(stdTx)
	//fmt.Println(txBytes)
	//fmt.Println(info.GetPubKey().Address())

	//fmt.Println(info.GetPubKey().Bytes())
	pass := info.GetPubKey().VerifyBytes(signByte1,sig)
	require.Equal(t, true, pass)
}
