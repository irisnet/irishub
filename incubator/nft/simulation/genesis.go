package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/modules/incubator/nft/internal/types"
)

const (
	kitties = "crypto-kitties"
	doggos  = "crypto-doggos"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	collections := types.NewCollections(types.NewCollection(doggos, types.NFTs{}), types.NewCollection(kitties, types.NFTs{}))
	var ownerships []types.Owner
	for _, acc := range simState.Accounts {
		// 10% of accounts own an NFT
		if simState.Rand.Intn(100) < 10 {
			baseNFT := types.NewBaseNFT(
				simulation.RandStringOfLength(simState.Rand, 10), // id
				acc.Address,
				simulation.RandStringOfLength(simState.Rand, 45), // tokenURI
			)

			var (
				idCollection types.IDCollection
				err          error
			)

			// 50% doggos and 50% kitties
			if simState.Rand.Intn(100) < 50 {
				collections[0], err = collections[0].AddNFT(&baseNFT)
				if err != nil {
					panic(err)
				}
				idCollection = types.NewIDCollection(doggos, []string{baseNFT.ID})
			} else {
				collections[1], err = collections[1].AddNFT(&baseNFT)
				if err != nil {
					panic(err)
				}
				idCollection = types.NewIDCollection(kitties, []string{baseNFT.ID})
			}

			ownership := types.NewOwner(acc.Address, idCollection)
			ownerships = append(ownerships, ownership)
		}
	}

	nftGenesis := types.NewGenesisState(ownerships, collections)

	fmt.Printf("Selected randomly generated NFT genesis state:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, nftGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(nftGenesis)
}
