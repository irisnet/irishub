package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"mods.irisnet.org/modules/nft/types"
)

// ---------------------------------------- Msgs --------------------------------------------------

func TestMsgTransferNFTValidateBasicMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, "", id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err := newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, uriHash, tokenData, "", address2.String())
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, uriHash, tokenData, address.String(), "")
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgTransferNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferNFTGetSignBytesMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	sortedBytes := newMsgTransferNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgTransferNFT","value":{"data":"{\"key1\":\"value1\",\"key2\":\"value2\"}","denom_id":"denom","id":"denom","name":"id1","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgTransferNFTGetSignersMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, uriHash, tokenData, address.String(), address2.String())
	signers := newMsgTransferNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestMsgEditNFTValidateBasicMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, uriHash, tokenData, "")

	err := newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT("", denom, nftName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, "", nftName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, denom, nftName, tokenURI, uriHash, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditNFTGetSignBytesMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, uriHash, tokenData, address.String())
	sortedBytes := newMsgEditNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgEditNFT","value":{"data":"{\"key1\":\"value1\",\"key2\":\"value2\"}","denom_id":"denom","id":"id1","name":"report","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgEditNFTGetSignersMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, uriHash, tokenData, address.String())
	signers := newMsgEditNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestMsgMsgMintNFTValidateBasicMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, uriHash, tokenData, "", address2.String())
	err := newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT("", denom, nftName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, "", nftName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, denom, nftName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintNFTGetSignBytesMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, uriHash, tokenData, address.String(), address2.String())
	sortedBytes := newMsgMintNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgMintNFT","value":{"data":"{\"key1\":\"value1\",\"key2\":\"value2\"}","denom_id":"denom","id":"id1","name":"report","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json","uri_hash":"uriHash"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgMsgBurnNFTValidateBasicMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT("", id, denom)
	err := newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address.String(), "", denom)
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address.String(), id, "")
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address.String(), id, denom)
	err = newMsgBurnNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgBurnNFTGetSignBytesMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(address.String(), id, denom)
	sortedBytes := newMsgBurnNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgBurnNFT","value":{"denom_id":"denom","id":"id1","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgBurnNFTGetSignersMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(address.String(), id, denom)
	signers := newMsgBurnNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}
