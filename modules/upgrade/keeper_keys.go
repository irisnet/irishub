package upgrade

import (
	"fmt"
	"strconv"
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	currentProposalId   = []byte("c/proposalId")
	currentVersionKey	= []byte("c/version")
	versionIDKey 		= "v/%s"		// v/<versionId>
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

func GetVersionIDKey(versionID int64) []byte {
	return []byte(fmt.Sprintf(versionIDKey, ToHexString(versionID)))
}

func GetProposalIDKey(proposalID int64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, ToHexString(proposalID)))
}

func GetStartHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf(startHeightKey, ToHexString(height)))
}

func GetSwitchKey(proposalID int64, switchVoterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf(switchKey, ToHexString(proposalID), switchVoterAddr.String()))
}

func ToHexString(i int64) string {
	hex := strconv.FormatInt(i, 16)
	var stringBuild bytes.Buffer
	for i:=0 ;i < 16 - len(hex); i++ {
		stringBuild.Write([]byte("0"))
	}
	stringBuild.Write([]byte(hex))
	return stringBuild.String()
}