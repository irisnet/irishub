package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/irisnet/irishub/v2/app/params"
)

const (
	testnetFile = "testnet-genesis-file"
	mainnetFile = "mainnet-genesis-file"
	outputFile  = "output-genesis-file"
)

// mergeGenesisCmd registers a sub-tree of commands to interact with
// local private key storage.
func mergeGenesisCmd(encodingConfig params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-genesis",
		Short: "merge genesis with testnet and mainnet",
		RunE: func(cmd *cobra.Command, _ []string) error {
			testnetGenesisPath, err := cmd.Flags().GetString(testnetFile)
			if err != nil {
				return err
			}

			mainnetGenesisPath, err := cmd.Flags().GetString(mainnetFile)
			if err != nil {
				return err
			}

			outputFile, err := cmd.Flags().GetString(outputFile)
			if err != nil {
				return err
			}

			testnetGenesis, err := genesisDocFromFile(testnetGenesisPath)
			if err != nil {
				return err
			}

			mainnetGenesis, err := genesisDocFromFile(mainnetGenesisPath)
			if err != nil {
				return err
			}

			return merge(encodingConfig.Marshaler, testnetGenesis, mainnetGenesis, outputFile)
		},
	}
	cmd.Flags().String(testnetFile, "", "irishub testnet genesis")
	cmd.Flags().String(mainnetFile, "", "irishub mainnet genesis")
	cmd.Flags().String(outputFile, "", "merged genesis")
	return cmd
}

func merge(cdc codec.Codec, testnet, mainnet *types.GenesisDoc, output string) (err error) {
	var mainnetAppState, testnetAppState map[string]json.RawMessage
	if err = tmjson.Unmarshal(mainnet.AppState, &mainnetAppState); err != nil {
		return err
	}

	if err = tmjson.Unmarshal(testnet.AppState, &testnetAppState); err != nil {
		panic(err)
	}
	mainnet.Validators = nil
	mainnetAppState["staking"] = testnetAppState["staking"]
	mainnetAppState["slashing"] = testnetAppState["slashing"]
	mainnetAppState["mint"] = testnetAppState["mint"]
	mainnetAppState["distribution"] = testnetAppState["distribution"]
	mainnetAppState["genutil"] = testnetAppState["genutil"]
	mainnetAppState["htlc"] = testnetAppState["htlc"]

	mergeBank(cdc, testnetAppState, mainnetAppState)
	mergeAuth(cdc, testnetAppState, mainnetAppState)
	mergeGov(cdc, testnetAppState, mainnetAppState)

	mainnet.InitialHeight = 0
	mainnet.ChainID = testnet.ChainID
	mainnet.AppState, err = tmjson.Marshal(mainnetAppState)
	if err != nil {
		return err
	}
	return mainnet.SaveAs(output)
}

var filterAddrs = map[string]bool{
	//distribution
	"iaa1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8jaydtw": true,
	//not_bonded_tokens_pool
	"iaa1tygms3xhhs3yv487phx3dw4a95jn7t7l5e40dj": true,
	//bonded_tokens_pool
	"iaa1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3qef7mx": true,
}

func mergeBank(cdc codec.Codec, testnet, mainnet map[string]json.RawMessage) {
	var bankState, testnetBankState banktypes.GenesisState
	cdc.MustUnmarshalJSON(mainnet["bank"], &bankState)

	//clean Supply
	bankState.Supply = sdk.NewCoins()

	//delete balance
	k := 0
	for _, balance := range bankState.Balances {
		if !filterAddrs[balance.Address] {
			bankState.Balances[k] = balance
			k++
		}
	}
	bankState.Balances = bankState.Balances[:k]
	//copy testnet balance to mainnet
	cdc.MustUnmarshalJSON(testnet["bank"], &testnetBankState)
	bankState.Balances = append(bankState.Balances, testnetBankState.Balances...)
	mainnet["bank"] = cdc.MustMarshalJSON(&bankState)
}

func mergeAuth(cdc codec.Codec, testnet, mainnet map[string]json.RawMessage) {
	var authState, testnetAuthState authtypes.GenesisState
	cdc.MustUnmarshalJSON(testnet["auth"], &testnetAuthState)
	cdc.MustUnmarshalJSON(mainnet["auth"], &authState)

	for _, account := range testnetAuthState.Accounts {
		authState.Accounts = append(authState.Accounts, account)
	}
	mainnet["auth"] = cdc.MustMarshalJSON(&authState)
}

func mergeGov(cdc codec.Codec, testnet, mainnet map[string]json.RawMessage) {
	var govState, testnetgovState govtypesv1.GenesisState
	cdc.MustUnmarshalJSON(testnet["gov"], &testnetgovState)
	cdc.MustUnmarshalJSON(mainnet["gov"], &govState)

	govState.DepositParams = testnetgovState.DepositParams
	govState.VotingParams = testnetgovState.VotingParams
	govState.TallyParams = testnetgovState.TallyParams
	mainnet["gov"] = cdc.MustMarshalJSON(&govState)
}

func genesisDocFromFile(genDocFile string) (*types.GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read GenesisDoc file: %w", err)
	}
	genDoc, err := types.GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading GenesisDoc at %s: %w", genDocFile, err)
	}
	return genDoc, nil
}
