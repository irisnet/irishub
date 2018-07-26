package upgrade

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

var (
	currentProposalId   = []byte("c/proposalId")
	currentVersionKey	= []byte("c/version")
	versionKey 			= "v/%s"		// v/<versionId>
	proposalIDKey 		= "p/%s"		// p/<proposalId>
	startHeightKey		= "h/%s"		// h/<height>
	switchKey			= "s/%s/%s"		// s/<proposalId>/<switchVoterAddress>
)

func GetCurrentProposalIdKey() []byte {
	return currentProposalId
}

func GetCurrentVersionKey() []byte {
	return currentVersionKey
}

func GetVersionIDKey(proposalID int64) []byte {
	return []byte(fmt.Sprintf(versionKey, strconv.FormatInt(proposalID, 16)))
}

func GetProposalIDKey(proposalID int64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, strconv.FormatInt(proposalID, 16)))
}

func GetStartHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf(startHeightKey, strconv.FormatInt(height, 16)))
}

func GetSwitchKey(proposalID int64, switchVoterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf(switchKey, strconv.FormatInt(proposalID, 16), switchVoterAddr.String()))
}

