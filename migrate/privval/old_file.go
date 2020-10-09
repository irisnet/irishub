package privval

import (
	"io/ioutil"
	"os"

	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/tempfile"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

// OldFilePV is the old version of the FilePV, pre v0.28.0.
type OldFilePV struct {
	Address       types.Address    `json:"address"`
	PubKey        crypto.PubKey    `json:"pub_key"`
	LastHeight    int64            `json:"last_height"`
	LastRound     int              `json:"last_round"`
	LastStep      int8             `json:"last_step"`
	LastSignature []byte           `json:"last_signature,omitempty"`
	LastSignBytes tmbytes.HexBytes `json:"last_signbytes,omitempty"`
	PrivKey       crypto.PrivKey   `json:"priv_key"`

	filePath string
}

// LoadOldFilePV loads an OldFilePV from the filePath.
func LoadOldFilePV(filePath string) (*OldFilePV, error) {
	pvJSONBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	pv := &OldFilePV{}
	err = cdc.UnmarshalJSON(pvJSONBytes, &pv)
	if err != nil {
		return nil, err
	}

	// overwrite pubkey and address for convenience
	pv.PubKey = pv.PrivKey.PubKey()
	pv.Address = pv.PubKey.Address()

	pv.filePath = filePath
	return pv, nil
}

// Upgrade convets the OldFilePV to the new FilePV, separating the immutable and mutable components,
// and persisting them to the keyFilePath and stateFilePath, respectively.
// It renames the original file by adding ".bak".
func (oldFilePV *OldFilePV) Upgrade(keyFilePath, stateFilePath string) *privval.FilePV {
	privKey := oldFilePV.PrivKey
	pvKey := privval.FilePVKey{
		PrivKey: privKey,
		PubKey:  privKey.PubKey(),
		Address: privKey.PubKey().Address(),
	}

	pvState := privval.FilePVLastSignState{
		Height:    oldFilePV.LastHeight,
		Round:     int32(oldFilePV.LastRound),
		Step:      oldFilePV.LastStep,
		Signature: oldFilePV.LastSignature,
		SignBytes: oldFilePV.LastSignBytes,
	}

	// Save the new PV files
	SavePvKey(pvKey, keyFilePath)
	SavePvState(pvState, stateFilePath)

	pv := &privval.FilePV{
		Key:           pvKey,
		LastSignState: pvState,
	}

	// Rename the old PV file
	err := os.Rename(oldFilePV.filePath, oldFilePV.filePath+".bak")
	if err != nil {
		panic(err)
	}
	return pv
}

func SavePvKey(pvKey privval.FilePVKey, keyFilePath string) {
	jsonBytes, err := tmjson.MarshalIndent(pvKey, "", "  ")
	if err != nil {
		panic(err)
	}
	err = tempfile.WriteFileAtomic(keyFilePath, jsonBytes, 0600)
	if err != nil {
		panic(err)
	}
}

func SavePvState(pvKey privval.FilePVLastSignState, stateFilePath string) {
	jsonBytes, err := tmjson.MarshalIndent(pvKey, "", "  ")
	if err != nil {
		panic(err)
	}
	err = tempfile.WriteFileAtomic(stateFilePath, jsonBytes, 0600)
	if err != nil {
		panic(err)
	}
}
