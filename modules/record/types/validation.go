package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateContents verifies whether the given contents are legal
func ValidateContents(contents ...Content) error {
	for i, content := range contents {
		if len(content.Digest) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "content[%d] digest missing", i)
		}
		if len(content.DigestAlgo) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "content[%d] digest algo missing", i)
		}
	}
	return nil
}
