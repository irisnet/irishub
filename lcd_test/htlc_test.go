//nolint:bodyclose
package lcdtest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/htlc"
	htlcrest "github.com/irisnet/irishub/modules/htlc/client/rest"
)

func TestHTLC(t *testing.T) {
	name := "sender"
	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20)))
	secret := htlc.HTLCSecret("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := htlc.HTLCHashLock(htlc.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock := uint64(50)

	kb, err := newKeybase()
	require.NoError(t, err)
	addrSender, _, err := CreateAddr(name, kb)
	addrTo, _, err := CreateAddr("to", kb)
	require.NoError(t, err)

	cleanup, _, _, port, err := InitializeLCD(1, []sdk.AccAddress{addrSender, addrTo}, true, []string{})
	require.NoError(t, err)
	defer cleanup()

	// create HTLC
	resultTx := createHTLC(
		t, port, name, kb,
		addrSender,
		addrTo,
		receiverOnOtherChain,
		amount,
		hashLock.String(),
		timeLock,
		timestamp,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	accSender := getAccount(t, port, addrSender)
	accTo := getAccount(t, port, addrTo)
	require.Equal(t, "99999975", accSender.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "100000000", accTo.GetCoins().AmountOf(sdk.DefaultBondDenom).String())

	// query HTLC
	resHTLC := queryHTLC(t, port, hashLock.String())
	require.Equal(t, addrSender.String(), resHTLC.Sender.String())
	require.Equal(t, addrTo.String(), resHTLC.To.String())
	require.Equal(t, receiverOnOtherChain, resHTLC.ReceiverOnOtherChain)
	require.Equal(t, amount, resHTLC.Amount)
	require.Equal(t, htlc.HTLCSecret{}, resHTLC.Secret)
	require.Equal(t, timestamp, resHTLC.Timestamp)
	require.Equal(t, htlc.OPEN, resHTLC.State)

	// wait for expired
	tests.WaitForHeight(resultTx.Height+55, port)

	// query HTLC
	resHTLC = queryHTLC(t, port, hashLock.String())
	require.Equal(t, addrSender.String(), resHTLC.Sender.String())
	require.Equal(t, addrTo.String(), resHTLC.To.String())
	require.Equal(t, receiverOnOtherChain, resHTLC.ReceiverOnOtherChain)
	require.Equal(t, amount, resHTLC.Amount)
	require.Equal(t, htlc.HTLCSecret{}, resHTLC.Secret)
	require.Equal(t, timestamp, resHTLC.Timestamp)
	require.Equal(t, htlc.EXPIRED, resHTLC.State)

	// refund HTLC
	resultTx = refundHTLC(
		t, port, name, kb,
		addrSender,
		hashLock.String(),
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	// query HTLC
	res, _ := Request(t, port, "GET", fmt.Sprintf("/htlc/htlcs/%s", hashLock), nil)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// create HTLC
	resultTx = createHTLC(
		t, port, name, kb,
		addrSender,
		addrTo,
		receiverOnOtherChain,
		amount,
		hashLock.String(),
		timeLock,
		timestamp,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	// query HTLC
	resHTLC = queryHTLC(t, port, hashLock.String())
	require.Equal(t, addrSender.String(), resHTLC.Sender.String())
	require.Equal(t, addrTo.String(), resHTLC.To.String())
	require.Equal(t, receiverOnOtherChain, resHTLC.ReceiverOnOtherChain)
	require.Equal(t, amount, resHTLC.Amount)
	require.Equal(t, htlc.HTLCSecret{}, resHTLC.Secret)
	require.Equal(t, timestamp, resHTLC.Timestamp)
	require.Equal(t, htlc.OPEN, resHTLC.State)

	// claim HTLC
	resultTx = claimHTLC(
		t, port, name, kb,
		addrSender,
		hashLock.String(),
		secret.String(),
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	accSender = getAccount(t, port, addrSender)
	accTo = getAccount(t, port, addrTo)
	require.Equal(t, "99999960", accSender.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "100000020", accTo.GetCoins().AmountOf(sdk.DefaultBondDenom).String())

	// query HTLC
	resHTLC = queryHTLC(t, port, hashLock.String())
	require.Equal(t, addrSender.String(), resHTLC.Sender.String())
	require.Equal(t, addrTo.String(), resHTLC.To.String())
	require.Equal(t, receiverOnOtherChain, resHTLC.ReceiverOnOtherChain)
	require.Equal(t, amount, resHTLC.Amount)
	require.Equal(t, secret, resHTLC.Secret)
	require.Equal(t, timestamp, resHTLC.Timestamp)
	require.Equal(t, htlc.COMPLETED, resHTLC.State)
}

// POST /htlc/htlcs
func createHTLC(
	t *testing.T,
	port string,
	name string,
	kb crkeys.Keybase,
	addrSender sdk.AccAddress,
	addrTo sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	hashLock string,
	timeLock uint64,
	timestamp uint64,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	createHTLCReq := htlcrest.CreateHTLCReq{
		BaseTx:               baseReq,
		Sender:               addrSender,
		To:                   addrTo,
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
		TimeLock:             timeLock,
		Timestamp:            timestamp,
	}

	req, err := cdc.MarshalJSON(createHTLCReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", "/htlc/htlcs", req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// POST /htlc/htlcs/{hash-lock}/claim
func claimHTLC(
	t *testing.T,
	port string,
	name string,
	kb crkeys.Keybase,
	addrSender sdk.AccAddress,
	hashLock string,
	secret string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	claimHTLCReq := htlcrest.ClaimHTLCReq{
		BaseTx: baseReq,
		Sender: addrSender,
		Secret: secret,
	}

	req, err := cdc.MarshalJSON(claimHTLCReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", fmt.Sprintf("/htlc/htlcs/%s/claim", hashLock), req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// POST /htlc/htlcs/{hash-lock}/refund
func refundHTLC(
	t *testing.T,
	port string,
	name string,
	kb crkeys.Keybase,
	addrSender sdk.AccAddress,
	hashLock string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	refundHTLCReq := htlcrest.RefundHTLCReq{
		BaseTx: baseReq,
		Sender: addrSender,
	}

	req, err := cdc.MarshalJSON(refundHTLCReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", fmt.Sprintf("/htlc/htlcs/%s/refund", hashLock), req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// GET /htlc/htlcs/{hash-lock}
func queryHTLC(t *testing.T, port string, hashLock string) (resHTLC htlc.HTLC) {
	res, body := Request(t, port, "GET", fmt.Sprintf("/htlc/htlcs/%s", hashLock), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var resp rest.ResponseWithHeight
	require.NoError(t, cdc.UnmarshalJSON([]byte(body), &resp))

	err := cdc.UnmarshalJSON(resp.Result, &resHTLC)
	require.NoError(t, err)

	return
}
