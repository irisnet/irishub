package rand

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr sdk.AccAddress
	testAddr  = sdk.AccAddress([]byte("testAddr"))
)

func TestNewMsgRequestRand(t *testing.T) {
	msg := NewMsgRequestRand(testAddr)

	require.Equal(t, testAddr, msg.Consumer)
}

func TestMsgRequestRandRoute(t *testing.T) {
	// build a MsgRequestRand
	msg := NewMsgRequestRand{
		Consumer: testAddr,
	}

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandValidation(t *testing.T) {
	testData := []struct {
		name       string
		consumer   sdk.AccAddress
		expectPass bool
	}{
		{"empty consumer", emptyAddr, false},
		{"valid consumer", testAddr, true},
	}

	for _, td := range testData {
		msg := NewMsgRequestRand(td.consumer)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"irishub/rand/MsgRequestRand\",\"value\":{\"consumer\":\"faa1damkuetjqqah8w\"}}"
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandGetSigners(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr)
	res := msg.GetSigners()

	expected := "[6F776E6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
