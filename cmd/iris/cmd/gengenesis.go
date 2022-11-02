package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/app"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/types"
)

const (
	testnetFile = "testnet-genesis-file"
	mainnetFile = "mainnet-genesis-file"
	outputFile  = "output-genesis-file"
)

// MergeGenesisCommands registers a sub-tree of commands to interact with
// local private key storage.
func MergeGenesisCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-genesis",
		Short: "merge genesis with testnet and mainnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			encodingConfig := app.MakeEncodingConfig()
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
	// mainnet.Validators = nil
	// mainnetAppState["staking"] = testnetAppState["staking"]
	// mainnetAppState["slashing"] = testnetAppState["slashing"]
	// mainnetAppState["mint"] = testnetAppState["mint"]
	// mainnetAppState["distribution"] = testnetAppState["distribution"]
	// mainnetAppState["genutil"] = testnetAppState["genutil"]
	// mainnetAppState["farm"] = testnetAppState["farm"]
	testnetAppState["nft"] = mainnetAppState["nft"]

	testnet.AppState, err = tmjson.Marshal(testnetAppState)
	if err != nil {
		return err
	}
	return testnet.SaveAs(output)
}

func mergeBank(cdc codec.Codec, testnet, mainnet map[string]json.RawMessage) {
	//distribution iaa1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8jaydtw
	//not_bonded_tokens_pool iaa1tygms3xhhs3yv487phx3dw4a95jn7t7l5e40dj
	// bonded_tokens_pool iaa1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3qef7mx
	// farm iaa1er8hq8es45ga8m580h8dp4m54vk6j9vckm4t8j
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
