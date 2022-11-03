package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/irisnet/irismod/modules/nft/exported"
	"github.com/irisnet/irismod/modules/nft/types"
)

// SaveNFT mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) SaveNFT(ctx sdk.Context, denomID,
	tokenID,
	tokenNm,
	tokenURI,
	tokenUriHash,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	nftMetadata := &types.NFTMetadata{
		Name: tokenNm,
		Data: tokenData,
	}
	data, err := codectypes.NewAnyWithValue(nftMetadata)
	if err != nil {
		return err
	}
	return k.nk.Mint(ctx, nft.NFT{
		ClassId: denomID,
		Id:      tokenID,
		Uri:     tokenURI,
		UriHash: tokenUriHash,
		Data:    data,
	}, receiver)
}

// UpdateNFT updates an already existing NFT
func (k Keeper) UpdateNFT(ctx sdk.Context, denomID,
	tokenID,
	tokenNm,
	tokenURI,
	tokenURIHash,
	tokenData string,
	owner sdk.AccAddress,
) error {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return err
	}

	if denom.UpdateRestricted {
		// if true , nobody can update the NFT under this denom
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "nobody can update the NFT under this denom %s", denomID)
	}

	// just the owner of NFT can edit
	if err := k.Authorize(ctx, denomID, tokenID, owner); err != nil {
		return err
	}

	if !types.Modified(tokenURI) &&
		!types.Modified(tokenURIHash) &&
		!types.Modified(tokenNm) &&
		!types.Modified(tokenData) {
		return nil
	}

	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return sdkerrors.Wrapf(types.ErrUnknownNFT, "nft ID %s not exists", tokenID)
	}

	token.Uri = types.Modify(token.Uri, tokenURI)
	token.UriHash = types.Modify(token.UriHash, tokenURIHash)
	if types.Modified(tokenNm) || types.Modified(tokenData) {
		var nftMetadata types.NFTMetadata
		if err := k.cdc.Unmarshal(token.Data.GetValue(), &nftMetadata); err != nil {
			return err
		}

		nftMetadata.Name = types.Modify(nftMetadata.Name, tokenNm)
		nftMetadata.Data = types.Modify(nftMetadata.Data, tokenData)
		data, err := codectypes.NewAnyWithValue(&nftMetadata)
		if err != nil {
			return err
		}
		token.Data = data
	}
	return k.nk.Update(ctx, token)
}

// TransferOwnership transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferOwnership(ctx sdk.Context, denomID,
	tokenID,
	tokenNm,
	tokenURI,
	tokenURIHash,
	tokenData string,
	srcOwner,
	dstOwner sdk.AccAddress,
) error {
	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return sdkerrors.Wrapf(types.ErrInvalidTokenID, "nft ID %s not exists", tokenID)
	}

	if err := k.Authorize(ctx, denomID, tokenID, srcOwner); err != nil {
		return err
	}

	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return err
	}

	tokenChanged := types.Modified(tokenURI) || types.Modified(tokenURIHash)
	tokenMetadataChanged := types.Modified(tokenNm) || types.Modified(tokenData)

	if denom.UpdateRestricted && (tokenChanged || tokenMetadataChanged) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "It is restricted to update NFT under this denom %s", denom.Id)
	}

	if !tokenChanged && !tokenMetadataChanged {
		return k.nk.Transfer(ctx, denomID, tokenID, dstOwner)
	}

	token.Uri = types.Modify(token.Uri, tokenURI)
	token.UriHash = types.Modify(token.UriHash, tokenURIHash)
	if tokenMetadataChanged {
		var nftMetadata types.NFTMetadata
		if err := k.cdc.Unmarshal(token.Data.GetValue(), &nftMetadata); err != nil {
			return err
		}

		nftMetadata.Name = types.Modify(nftMetadata.Name, tokenNm)
		nftMetadata.Data = types.Modify(nftMetadata.Data, tokenData)
		data, err := codectypes.NewAnyWithValue(&nftMetadata)
		if err != nil {
			return err
		}
		token.Data = data
	}

	if err := k.nk.Update(ctx, token); err != nil {
		return err
	}
	return k.nk.Transfer(ctx, denomID, tokenID, dstOwner)
}

// RemoveNFT deletes a specified NFT
func (k Keeper) RemoveNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if err := k.Authorize(ctx, denomID, tokenID, owner); err != nil {
		return err
	}
	return k.nk.Burn(ctx, denomID, tokenID)
}

// GetNFT gets the specified NFT
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error) {
	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrUnknownNFT, "not found NFT: %s", denomID)
	}

	var nftMetadata types.NFTMetadata
	if err := k.cdc.Unmarshal(token.Data.GetValue(), &nftMetadata); err != nil {
		return nil, err
	}

	owner := k.nk.GetOwner(ctx, denomID, tokenID)
	return types.BaseNFT{
		Id:      token.Id,
		Name:    nftMetadata.Name,
		URI:     token.Uri,
		Data:    nftMetadata.Data,
		Owner:   owner.String(),
		UriHash: token.UriHash,
	}, nil
}

// GetNFTs returns all NFTs by the specified denom ID
func (k Keeper) GetNFTs(ctx sdk.Context, denom string) (nfts []exported.NFT, err error) {
	tokens := k.nk.GetNFTsOfClass(ctx, denom)
	for _, token := range tokens {
		var nftMetadata types.NFTMetadata
		if err := k.cdc.Unmarshal(token.Data.GetValue(), &nftMetadata); err != nil {
			return nil, err
		}
		nfts = append(nfts, types.BaseNFT{
			Id:      token.GetId(),
			Name:    nftMetadata.Name,
			URI:     token.GetUri(),
			UriHash: token.GetUriHash(),
			Data:    nftMetadata.Data,
			Owner:   k.nk.GetOwner(ctx, denom, token.GetId()).String(),
		})
	}
	return nfts, nil
}

// Authorize checks if the sender is the owner of the given NFT
// Return the NFT if true, an error otherwise
func (k Keeper) Authorize(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if !owner.Equals(k.nk.GetOwner(ctx, denomID, tokenID)) {
		return sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}
	return nil
}

// HasNFT checks if the specified NFT exists
func (k Keeper) HasNFT(ctx sdk.Context, denomID, tokenID string) bool {
	return k.nk.HasNFT(ctx, denomID, tokenID)
}
