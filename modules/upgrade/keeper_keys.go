package upgrade

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	CurrentVersionIDKey	= []byte("c/")
	VersionKey 			= "v/%010d/"		// v/<proposalId>
	SwitchKey			= "s/%010d/%d/"		// s/<proposalId>/<switchVoterAddress>
	VersionListKey      = []byte("l/")		// list of the version_key ordered by proposalId
)

func GetCurrentVersionIDKey() []byte {
	return CurrentVersionIDKey
}

func GetVersionKey(proposalID int64) []byte {
	return []byte(fmt.Sprintf(VersionKey, proposalID))
}

func GetSwitchKey(proposalID int64, switchVoterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf(SwitchKey, proposalID, switchVoterAddr))
}

func GetVersionListKey() []byte {
	return VersionListKey
}
