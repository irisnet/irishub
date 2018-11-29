package upgrade

import (
	"fmt"
	"strconv"
	"bytes"
	sdk "github.com/irisnet/irishub/types"
)

var (
	currentVersionKey	        = []byte("c/version")
	versionIDKey 		        = "v/%s"		// v/<versionId>
	proposalIDKey 		        = "p/%s"		// p/<proposalId>
	startHeightKey		        = "h/%s"		// h/<height>
	switchKey			        = "s/%s/%s"		// s/<proposalId>/<switchVoterAddress>
	DoingSwitchKey				= []byte("d")		// whether system is doing switch
)

func GetCurrentVersionKey() []byte {
	return currentVersionKey
}

func GetVersionIDKey(versionID int64) []byte {
	return []byte(fmt.Sprintf(versionIDKey, IntToHexString(versionID)))
}

func GetProposalIDKey(proposalID uint64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, UintToHexString(proposalID)))
}

func GetStartHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf(startHeightKey, IntToHexString(height)))
}

func GetSwitchKey(proposalID uint64, switchVoterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf(switchKey, UintToHexString(proposalID), switchVoterAddr.String()))
}

func IntToHexString(i int64) string {
	hex := strconv.FormatInt(i, 16)
	var stringBuild bytes.Buffer
	for i:=0 ;i < 16 - len(hex); i++ {
		stringBuild.Write([]byte("0"))
	}
	stringBuild.Write([]byte(hex))
	return stringBuild.String()
}
func UintToHexString(i uint64) string {
	hex := strconv.FormatUint(i, 16)
	var stringBuild bytes.Buffer
	for i:=0 ;i < 16 - len(hex); i++ {
		stringBuild.Write([]byte("0"))
	}
	stringBuild.Write([]byte(hex))
	return stringBuild.String()
}
func GetDoingSwitchKey() []byte {
	return DoingSwitchKey
}
