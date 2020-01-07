package lcdtest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/mint"
)

func TestMint(t *testing.T) {
	name := "sender"
	kb, err := keys.NewKeyringFromDir(InitClientHome(""), nil)
	require.NoError(t, err)
	addr, _, err := CreateAddr(name, kb)
	require.NoError(t, err)

	cleanup, _, _, port, err := InitializeLCD(1, []sdk.AccAddress{addr}, true)
	require.NoError(t, err)
	defer cleanup()

	// query params
	params := queryParams(t, port)
	err = params.Validate()
	require.NoError(t, err)
}

// GET /mint/parameters
func queryParams(t *testing.T, port string) mint.Params {
	res, body := Request(t, port, "GET", "/mint/parameters", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var params mint.Params
	require.NoError(t, cdc.UnmarshalJSON(extractResultFromResponse(t, []byte(body)), &params))

	return params
}
