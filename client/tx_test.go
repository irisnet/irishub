package client

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

func TestMsgDelegate(t *testing.T) {
	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	coins := sdk.NewCoin("iris", 10)
	var delegate = stake.NewMsgDelegate(addr1,addr1,coins)
	var sig2 = delegate.GetSignBytes()
	fmt.Println(string(sig2))
	fmt.Println(sig2)
}


func makeDelegateTx() sdk.Msg{
	addr1, _ := sdk.AccAddressFromHex("DA16415D899555501471A3763F3C0B5EA380E287")
	coins := sdk.NewCoin("iris", 10)
	var delegate = stake.NewMsgDelegate(addr1,addr1,coins)
	return delegate
}

func makeSendTx() sdk.Msg{
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
	chainID := "fuxi-1001"
	accnum := int64(0)
	sequence := int64(9)
	fee := auth.NewStdFee(10000, sdk.NewCoin("iris", 0))
	memo := "test"

	msgs := []sdk.Msg{
		//makeSendTx(),
		makeDelegateTx(),
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

	pass := info.GetPubKey().VerifyBytes(signByte1,sig)
	require.Equal(t, true, pass)
}