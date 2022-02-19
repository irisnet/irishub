package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irismod/modules/mt/types"
)

// ---------------------------------------- Msgs --------------------------------------------------

func TestMsgTransferMTValidateBasicMethod(t *testing.T) {
	newMsgTransferMT := types.NewMsgTransferMT(denomID, "", id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err := newMsgTransferMT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferMT = types.NewMsgTransferMT(denomID, denom, "", tokenURI, uriHash, tokenData, "", address2.String())
	err = newMsgTransferMT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferMT = types.NewMsgTransferMT(denomID, denom, "", tokenURI, uriHash, tokenData, address.String(), "")
	err = newMsgTransferMT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferMT = types.NewMsgTransferMT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgTransferMT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferMTGetSignBytesMethod(t *testing.T) {
	newMsgTransferMT := types.NewMsgTransferMT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	sortedBytes := newMsgTransferMT.GetSignBytes()
	expected := `{"type":"irismod/mt/MsgTransferMT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"denom","name":"id1","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgTransferMTGetSignersMethod(t *testing.T) {
	newMsgTransferMT := types.NewMsgTransferMT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	signers := newMsgTransferMT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestMsgEditMTValidateBasicMethod(t *testing.T) {
	newMsgEditMT := types.NewMsgEditMT(id, denom, mtName, tokenURI, uriHash, tokenData, "")

	err := newMsgEditMT.ValidateBasic()
	require.Error(t, err)

	newMsgEditMT = types.NewMsgEditMT("", denom, mtName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditMT.ValidateBasic()
	require.Error(t, err)

	newMsgEditMT = types.NewMsgEditMT(id, "", mtName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditMT.ValidateBasic()
	require.Error(t, err)

	newMsgEditMT = types.NewMsgEditMT(id, denom, mtName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditMT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditMTGetSignBytesMethod(t *testing.T) {
	newMsgEditMT := types.NewMsgEditMT(id, denom, mtName, tokenURI, uriHash, tokenData, address.String())
	sortedBytes := newMsgEditMT.GetSignBytes()
	expected := `{"type":"irismod/mt/MsgEditMT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"id1","name":"report","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgEditMTGetSignersMethod(t *testing.T) {
	newMsgEditMT := types.NewMsgEditMT(id, denom, mtName, tokenURI, uriHash, tokenData, address.String())
	signers := newMsgEditMT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestMsgMsgMintMTValidateBasicMethod(t *testing.T) {
	newMsgMintMT := types.NewMsgMintMT(id, denom, mtName, tokenURI, uriHash, tokenData, "", address2.String())
	err := newMsgMintMT.ValidateBasic()
	require.Error(t, err)

	newMsgMintMT = types.NewMsgMintMT("", denom, mtName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintMT.ValidateBasic()
	require.Error(t, err)

	newMsgMintMT = types.NewMsgMintMT(id, "", mtName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintMT.ValidateBasic()
	require.Error(t, err)

	newMsgMintMT = types.NewMsgMintMT(id, denom, mtName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintMT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintMTGetSignBytesMethod(t *testing.T) {
	newMsgMintMT := types.NewMsgMintMT(id, denom, mtName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	sortedBytes := newMsgMintMT.GetSignBytes()
	expected := `{"type":"irismod/mt/MsgMintMT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"id1","name":"report","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgMsgBurnMTValidateBasicMethod(t *testing.T) {
	newMsgBurnMT := types.NewMsgBurnMT("", id, denom)
	err := newMsgBurnMT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnMT = types.NewMsgBurnMT(address.String(), "", denom)
	err = newMsgBurnMT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnMT = types.NewMsgBurnMT(address.String(), id, "")
	err = newMsgBurnMT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnMT = types.NewMsgBurnMT(address.String(), id, denom)
	err = newMsgBurnMT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgBurnMTGetSignBytesMethod(t *testing.T) {
	newMsgBurnMT := types.NewMsgBurnMT(address.String(), id, denom)
	sortedBytes := newMsgBurnMT.GetSignBytes()
	expected := `{"type":"irismod/mt/MsgBurnMT","value":{"denom_id":"denom","id":"id1","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgBurnMTGetSignersMethod(t *testing.T) {
	newMsgBurnMT := types.NewMsgBurnMT(address.String(), id, denom)
	signers := newMsgBurnMT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}
