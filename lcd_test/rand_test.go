//nolint:bodyclose
package lcdtest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/rand"
	randrest "github.com/irisnet/irishub/modules/rand/client/rest"
)

func TestRand(t *testing.T) {
	name := "sender"
	blockInterval := uint64(50)
	kb, err := newKeybase()
	require.NoError(t, err)
	addr, _, err := CreateAddr(name, kb)
	require.NoError(t, err)

	cleanup, _, _, port, err := InitializeLCD(1, []sdk.AccAddress{addr}, true, []string{})
	require.NoError(t, err)
	defer cleanup()

	// request rand
	resultTx := requestRand(t, port, name, kb, addr, blockInterval)
	requestID := resultTx.Logs[0].Events[1].Attributes[0].Value
	generateHeight, err := strconv.ParseInt(resultTx.Logs[0].Events[1].Attributes[1].Value, 10, 64)
	txHash := resultTx.TxHash
	require.NoError(t, err)
	tests.WaitForHeight(resultTx.Height+1, port)

	// query queue
	requests := queryQueue(t, port, generateHeight)
	require.Equal(t, strings.ToLower(txHash), hex.EncodeToString(requests[0].TxHash))
	require.Equal(t, resultTx.Height, requests[0].Height)
	tests.WaitForHeight(resultTx.Height+55, port)

	// query rand
	readableRand := queryRand(t, port, requestID)
	require.Equal(t, generateHeight, readableRand.Height)
}

// POST /rand/rands
func requestRand(
	t *testing.T, port string, name string, kb crkeys.Keybase,
	addrSender sdk.AccAddress, blockInterval uint64,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	requestRandReq := randrest.RequestRandReq{
		BaseReq:       baseReq,
		Consumer:      addrSender,
		BlockInterval: blockInterval,
	}

	req, err := cdc.MarshalJSON(requestRandReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", "/rand/rands", req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// GET /rand/rands/{request-id}
func queryRand(t *testing.T, port string, requestID string) (readableRand rand.ReadableRand) {
	res, body := Request(t, port, "GET", fmt.Sprintf("/rand/rands/%s", requestID), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var resp rest.ResponseWithHeight
	require.NoError(t, cdc.UnmarshalJSON([]byte(body), &resp))

	err := cdc.UnmarshalJSON(resp.Result, &readableRand)
	require.NoError(t, err)

	return
}

// GET /rand/queue
func queryQueue(t *testing.T, port string, height int64) (requests []rand.Request) {
	res, body := Request(t, port, "GET", fmt.Sprintf("/rand/queue?gen-height=%d", height), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var resp rest.ResponseWithHeight
	require.NoError(t, cdc.UnmarshalJSON([]byte(body), &resp))

	err := cdc.UnmarshalJSON(resp.Result, &requests)
	require.NoError(t, err)

	return
}
