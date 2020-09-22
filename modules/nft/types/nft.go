package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/exported"
)

var _ exported.NFT = BaseNFT{}

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id, name string, owner sdk.AccAddress, tokenURI, tokenData string) BaseNFT {
	return BaseNFT{
		Id:    strings.ToLower(strings.TrimSpace(id)),
		Name:  strings.TrimSpace(name),
		Owner: owner,
		URI:   strings.TrimSpace(tokenURI),
		Data:  strings.TrimSpace(tokenData),
	}
}

func (bnft BaseNFT) GetID() string {
	return bnft.Id
}

func (bnft BaseNFT) GetName() string {
	return bnft.Id
}

func (bnft BaseNFT) GetOwner() sdk.AccAddress {
	return bnft.Owner
}

func (bnft BaseNFT) GetURI() string {
	return bnft.URI
}

func (bnft BaseNFT) GetData() string {
	return bnft.Data
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts)
}

func ValidateTokenID(tokenID string) error {
	tokenID = strings.TrimSpace(tokenID)
	if len(tokenID) < MinDenomLen || len(tokenID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "invalid tokenID %s, only accepts value [%d, %d]", tokenID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(tokenID) || !IsAlphaNumeric(tokenID) {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "invalid tokenID %s, only accepts alphanumeric characters,and begin with an english letter", tokenID)
	}
	return nil
}

func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return sdkerrors.Wrapf(ErrInvalidTokenURI, "invalid tokenURI %s, only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}
