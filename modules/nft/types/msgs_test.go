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
		tokenURI, tokenData, address, address2,
	)
	require.Equal(t, newMsgTransferNFT.Sender, address)
	require.Equal(t, newMsgTransferNFT.Recipient, address2)
	require.Equal(t, newMsgTransferNFT.Denom, denom)
	require.Equal(t, newMsgTransferNFT.Id, denomID)
}

func TestMsgTransferNFTValidateBasicMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, "", id, tokenURI, tokenData, address, address2)
	err := newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, tokenData, nil, address2)
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, "", tokenURI, tokenData, address, nil)
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address, address2)
	err = newMsgTransferNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferNFTGetSignBytesMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address, address2)
	sortedBytes := newMsgTransferNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), `{"type":"irismod/nft/MsgTransferNFT","value":{"data":"https://google.com/token-1.json","denom":"denom","id":"denom","name":"id1","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`)
}

func TestMsgTransferNFTGetSignersMethod(t *testing.T) {
	newMsgTransferNFT := types.NewMsgTransferNFT(denomID, denom, id, tokenURI, tokenData, address, address2)
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
		tokenData, address,
	)

	require.Equal(t, newMsgEditNFT.Sender.String(), address.String())
	require.Equal(t, newMsgEditNFT.Id, id)
	require.Equal(t, newMsgEditNFT.Denom, denom)
	require.Equal(t, newMsgEditNFT.URI, tokenURI)
}

func TestMsgEditNFTValidateBasicMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, nil)

	err := newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT("", denom, nftName, tokenURI, tokenData, address)
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, "", nftName, tokenURI, tokenData, address)
	err = newMsgEditNFT.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFT = types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address)
	err = newMsgEditNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditNFTGetSignBytesMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address)
	sortedBytes := newMsgEditNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), `{"type":"irismod/nft/MsgEditNFT","value":{"data":"https://google.com/token-1.json","denom":"denom","id":"id1","name":"report","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`)
}

func TestMsgEditNFTGetSignersMethod(t *testing.T) {
	newMsgEditNFT := types.NewMsgEditNFT(id, denom, nftName, tokenURI, tokenData, address)
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
		tokenData, address, address2,
	)

	require.Equal(t, newMsgMintNFT.Sender.String(), address.String())
	require.Equal(t, newMsgMintNFT.Recipient.String(), address2.String())
	require.Equal(t, newMsgMintNFT.Id, id)
	require.Equal(t, newMsgMintNFT.Denom, denom)
	require.Equal(t, newMsgMintNFT.URI, tokenURI)
}

func TestMsgMsgMintNFTValidateBasicMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, nil, address2)
	err := newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT("", denom, nftName, tokenURI, tokenData, address, address2)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, "", nftName, tokenURI, tokenData, address, address2)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, address, address2)
	err = newMsgMintNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintNFTGetSignBytesMethod(t *testing.T) {
	newMsgMintNFT := types.NewMsgMintNFT(id, denom, nftName, tokenURI, tokenData, address, address2)
	sortedBytes := newMsgMintNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), `{"type":"irismod/nft/MsgMintNFT","value":{"data":"https://google.com/token-1.json","denom":"denom","id":"id1","name":"report","recipient":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","uri":"https://google.com/token-1.json"}}`)
}

func TestNewMsgBurnNFT(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(
		address,
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
	)

	require.Equal(t, newMsgBurnNFT.Sender.String(), address.String())
	require.Equal(t, newMsgBurnNFT.Id, id)
	require.Equal(t, newMsgBurnNFT.Denom, denom)
}

func TestMsgMsgBurnNFTValidateBasicMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(nil, id, denom)
	err := newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address, "", denom)
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address, id, "")
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = types.NewMsgBurnNFT(address, id, denom)
	err = newMsgBurnNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgBurnNFTGetSignBytesMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(address, id, denom)
	sortedBytes := newMsgBurnNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), `{"type":"irismod/nft/MsgBurnNFT","value":{"denom":"denom","id":"id1","sender":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq"}}`)
}

func TestMsgBurnNFTGetSignersMethod(t *testing.T) {
	newMsgBurnNFT := types.NewMsgBurnNFT(address, id, denom)
	signers := newMsgBurnNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}
