package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	mt "github.com/irisnet/irismod/mt/types"
)

const (
	prefixDenomID   = "denom:id:"
	prefixDenomName = "denom:name:"
	prefixMtID      = "mt:id"

	lenMTs = 10         // MTs number under a denom
	supply = uint64(10) // supply for each MT
)

// genMTs returns MTs for a denom
func genMTs(r *rand.Rand) []mt.MT {
	mts := make([]mt.MT, lenMTs)
	for i := 0; i < lenMTs; i++ {
		mts[i] = mt.MT{
			Id:     prefixMtID + simtypes.RandStringOfLength(r, 10),
			Supply: supply,
			Data:   []byte(simtypes.RandStringOfLength(r, 10)),
		}
	}
	return mts
}

// genCollections returns a slice of mt collection
func genCollections(r *rand.Rand, accounts []simtypes.Account) []mt.Collection {
	collections := make([]mt.Collection, len(accounts))
	for i := 0; i < len(accounts); i++ {
		collections[i] = mt.Collection{
			Denom: &mt.Denom{
				Id:    prefixDenomID + simtypes.RandStringOfLength(r, 10),
				Name:  prefixDenomName + simtypes.RandStringOfLength(r, 10),
				Data:  []byte(simtypes.RandStringOfLength(r, 10)),
				Owner: accounts[i].Address.String(),
			},
			Mts: genMTs(r),
		}
	}
	return collections
}

// genDenomBalances generates DenomBalances for each account.
// mts must belong to denomId.
func genDenomBalances(r *rand.Rand, denomId string, mts []mt.MT, accounts []simtypes.Account) []mt.DenomBalance {
	denomBalances := make([]mt.DenomBalance, len(accounts))
	for i := 0; i < len(accounts); i++ {
		balances := make([]mt.Balance, len(mts))

		// amount evenly distributed and sum-up not exceeding the total supply
		for j := 0; j < len(mts); j++ {
			balances[j] = mt.Balance{
				MtId:   mts[j].Id,
				Amount: mts[j].Supply / uint64(len(accounts)),
			}
		}

		denomBalances[i] = mt.DenomBalance{
			DenomId:  denomId,
			Balances: balances,
		}
	}

	return denomBalances
}

// genOwners returns a slice of mt owner
func genOwners(r *rand.Rand, collections []mt.Collection, accounts []simtypes.Account) []mt.Owner {
	owners := make([]mt.Owner, len(accounts))
	for i := 0; i < len(accounts); i++ {
		collection := collections[i]
		owners[i] = mt.Owner{
			Address: accounts[i].Address.String(),
			Denoms:  genDenomBalances(r, collection.Denom.Id, collection.Mts, accounts),
		}
	}

	return owners
}

// RandomizedGenState generates a random GenesisState for mt.
func RandomizedGenState(simState *module.SimulationState) {
	var (
		collections []mt.Collection
		owners      []mt.Owner
		accLen      int = 10
	)

	if len(simState.Accounts) < accLen {
		accLen = len(simState.Accounts)
	}

	accs := simState.Accounts[:accLen]
	simState.AppParams.GetOrGenerate(simState.Cdc, "mt", &collections, simState.Rand,
		func(r *rand.Rand) {
			collections = genCollections(r, accs)
		},
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, "mt", &owners, simState.Rand,
		func(r *rand.Rand) {
			owners = genOwners(r, collections, accs)
		},
	)

	mtGenesis := &mt.GenesisState{
		Collections: collections,
		Owners:      owners,
	}

	bz, err := json.MarshalIndent(mtGenesis, "", " ")
	if err != nil {
		fmt.Printf("Selected randomly generated %s parameters:\n%s\n", mt.ModuleName, bz)
	}

	simState.GenState[mt.ModuleName] = simState.Cdc.MustMarshalJSON(mtGenesis)

}
