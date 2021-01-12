package v0_16

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/irisnet/irishub/migrate/v0_16/auth"
	"github.com/irisnet/irishub/migrate/v0_16/coinswap"
	"github.com/irisnet/irishub/migrate/v0_16/distribution"
	"github.com/irisnet/irishub/migrate/v0_16/gov"
	"github.com/irisnet/irishub/migrate/v0_16/guardian"
	"github.com/irisnet/irishub/migrate/v0_16/htlc"
	"github.com/irisnet/irishub/migrate/v0_16/mint"
	"github.com/irisnet/irishub/migrate/v0_16/rand"
	"github.com/irisnet/irishub/migrate/v0_16/service"
	"github.com/irisnet/irishub/migrate/v0_16/slashing"
	"github.com/irisnet/irishub/migrate/v0_16/stake"
	"github.com/irisnet/irishub/migrate/v0_16/upgrade"
)

type GenesisFileState struct {
	Accounts     []GenesisFileAccount      `json:"accounts"`
	AuthData     auth.GenesisState         `json:"auth"`
	StakeData    stake.GenesisState        `json:"stake"`
	MintData     mint.GenesisState         `json:"mint"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	UpgradeData  upgrade.GenesisState      `json:"upgrade"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	ServiceData  service.GenesisState      `json:"service"`
	GuardianData guardian.GenesisState     `json:"guardian"`
	RandData     rand.GenesisState         `json:"rand"`
	SwapData     coinswap.GenesisState     `json:"swap"`
	HtlcData     htlc.GenesisState         `json:"htlc"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
}

type GenesisFileAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

// GenesisDoc defines the initial conditions for a tendermint blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime     time.Time                  `json:"genesis_time"`
	ChainID         string                     `json:"chain_id"`
	ConsensusParams ConsensusParams            `json:"consensus_params,omitempty"`
	Validators      []tmtypes.GenesisValidator `json:"validators,omitempty"`
	AppHash         tmbytes.HexBytes           `json:"app_hash"`
	AppState        json.RawMessage            `json:"app_state,omitempty"`
}

// ConsensusParams contains consensus critical parameters that determine the
// validity of blocks.
type ConsensusParams struct {
	BlockSize abcitypes.BlockParams   `json:"block_size"`
	Evidence  EvidenceParams          `json:"evidence"`
	Validator tmproto.ValidatorParams `json:"validator"`
}

// EvidenceParams determine how we handle evidence of malfeasance
type EvidenceParams struct {
	MaxAge int64 `json:"max_age"` // only accept new evidence more recent than this
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (*GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error reading GenesisDoc at %v", genDocFile))
	}
	return genDoc, nil
}

func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := tmjson.Unmarshal(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// ValidateAndComplete checks that all necessary fields are present
// and fills in defaults for optional fields left empty
func (genDoc *GenesisDoc) ValidateAndComplete() error {
	if genDoc.ChainID == "" {
		return errors.New("Genesis doc must include non-empty chain_id")
	}
	if len(genDoc.ChainID) > tmtypes.MaxChainIDLen {
		return errors.Errorf("chain_id in genesis doc is too long (max: %d)", tmtypes.MaxChainIDLen)
	}

	if err := genDoc.ConsensusParams.Validate(); err != nil {
		return err
	}

	for i, v := range genDoc.Validators {
		if v.Power == 0 {
			return errors.Errorf("The genesis file cannot contain validators with no voting power: %v", v)
		}
		if len(v.Address) > 0 && !bytes.Equal(v.PubKey.Address(), v.Address) {
			return errors.Errorf("Incorrect address for validator %v in the genesis file, should be %v", v, v.PubKey.Address())
		}
		if len(v.Address) == 0 {
			genDoc.Validators[i].Address = v.PubKey.Address()
		}
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	return nil
}

// Validate validates the ConsensusParams to ensure all values are within their
// allowed limits, and returns an error if they are not.
func (params ConsensusParams) Validate() error {
	if params.BlockSize.MaxBytes <= 0 {
		return errors.Errorf("BlockSize.MaxBytes must be greater than 0. Got %d",
			params.BlockSize.MaxBytes)
	}
	if params.BlockSize.MaxBytes > tmtypes.MaxBlockSizeBytes {
		return errors.Errorf("BlockSize.MaxBytes is too big. %d > %d",
			params.BlockSize.MaxBytes, tmtypes.MaxBlockSizeBytes)
	}

	if params.BlockSize.MaxGas < -1 {
		return errors.Errorf("BlockSize.MaxGas must be greater or equal to -1. Got %d",
			params.BlockSize.MaxGas)
	}

	if params.Evidence.MaxAge <= 0 {
		return errors.Errorf("EvidenceParams.MaxAge must be greater than 0. Got %d",
			params.Evidence.MaxAge)
	}

	if len(params.Validator.PubKeyTypes) == 0 {
		return errors.New("len(Validator.PubKeyTypes) must be greater than 0")
	}

	// Check if keyType is a known ABCIPubKeyType
	for i := 0; i < len(params.Validator.PubKeyTypes); i++ {
		keyType := params.Validator.PubKeyTypes[i]
		if _, ok := tmtypes.ABCIPubKeyTypesToNames[keyType]; !ok {
			return errors.Errorf("params.Validator.PubKeyTypes[%d], %s, is an unknown pubkey type",
				i, keyType)
		}
	}

	return nil
}
