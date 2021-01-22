package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irismod/modules/nft/types"
)

// ---------------------------------------- Msgs ---------------------------------------------------

func TestNewMsgTransferNFT(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(
		fmt.Sprintf("     %s     ", denomID),
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", id),
		tokenURI, tokenData, address.String(), address2.String(),
	)
	msgTransferNFT := newMsgTransferNFT.Normalize()
	require.Equal(t, msgTransferNFT.Sender, address.String())
	require.Equal(t, msgTransferNFT.Recipient, address2.String())
	require.Equal(t, msgTransferNFT.DenomId, denom)
	require.Equal(t, msgTransferNFT.Id, denomID)
}

func TestMsgTransferNFTValidateBasicMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, "", id, tokenURI, tokenData, address.String(), address2.String())
	err := newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, tokenData, "", address2.String())
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, tokenData, address.String(), "")
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address.String(), address2.String())
	err = newMsgTransferNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferNFTGetSignBytesMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address.String(), address2.String())
	sortedBytes := newMsgTransferNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgTransferNFT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"denom","name":"id1","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgTransferNFTGetSignersMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address.String(), address2.String())
	signers := newMsgTransferNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestNewMsgEditNFT(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", nftName),
		fmt.Sprintf("     %s     ", tokenURI),
		tokenData, address.String(),
	)
	msgEditNFT := newMsgEditNFT.Normalize()
	require.Equal(t, msgEditNFT.Sender, address.String())
	require.Equal(t, msgEditNFT.Id, id)
	require.Equal(t, msgEditNFT.DenomId, denom)
	require.Equal(t, msgEditNFT.URI, tokenURI)
}

func TestMsgEditNFTValidateBasicMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, "")

	err := newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT("", denom, nftName, tokenURI, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, "", nftName, tokenURI, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address.String())
	err = newMsgEditNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditNFTGetSignBytesMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address.String())
	sortedBytes := newMsgEditNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgEditNFT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"id1","name":"report","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestMsgEditNFTGetSignersMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address.String())
	signers := newMsgEditNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestNewMsgMintNFT(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", nftName),
		fmt.Sprintf("     %s     ", tokenURI),
		tokenData, address.String(), address2.String(),
	)
	msgMintNFT := newMsgMintNFT.Normalize()
	require.Equal(t, msgMintNFT.Sender, address.String())
	require.Equal(t, msgMintNFT.Recipient, address2.String())
	require.Equal(t, msgMintNFT.Id, id)
	require.Equal(t, msgMintNFT.DenomId, denom)
	require.Equal(t, msgMintNFT.URI, tokenURI)
}

func TestMsgMsgMintNFTValidateBasicMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, "", address2.String())
	err := newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT("", denom, nftName, tokenURI, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, "", nftName, tokenURI, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, address.String(), address2.String())
	err = newMsgMintNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintNFTGetSignBytesMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, address.String(), address2.String())
	sortedBytes := newMsgMintNFT.GetSignBytes()
	expected := `{"type":"irismod/nft/MsgMintNFT","value":{"data":"https://google.com/token-1.json","denom_id":"denom","id":"id1","name":"report","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`
	require.Equal(t, expected, string(sortedBytes))
}

func TestNewMsgBurnNFT(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(
		address.String(),
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
	)
	msgBurnNFT := newMsgBurnNFT.Normalize()
	require.Equal(t, msgBurnNFT.Sender, address.String())
	require.Equal(t, msgBurnNFT.Id, id)
	require.Equal(t, msgBurnNFT.DenomId, denom)
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
