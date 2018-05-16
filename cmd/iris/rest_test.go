package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"testing"
	"github.com/tendermint/tmlibs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/commands"
	rest "github.com/cosmos/cosmos-sdk/client/rest"
	byteTx "github.com/irisnet/iris-hub/rest"
	coinrest "github.com/cosmos/cosmos-sdk/modules/coin/rest"
	noncerest "github.com/cosmos/cosmos-sdk/modules/nonce/rest"
	rolerest "github.com/cosmos/cosmos-sdk/modules/roles/rest"

	stakerest "github.com/MrXJC/gaia/modules/stake/rest"
)

func TestRest(t *testing.T) {


	router := mux.NewRouter()

	rootDir := viper.GetString(cli.HomeFlag)
	keyMan := client.GetKeyManager(rootDir)
	serviceKeys := rest.NewServiceKeys(keyMan)
	serviceByteTx := byteTx.NewServiceByteTx(keyMan)
	serviceTxs := rest.NewServiceTxs(commands.GetNode())


	routeRegistrars := []func(*mux.Router) error{
		// rest.Keys handlers
		serviceKeys.RegisterCRUD,

		// Coin handlers (Send, Query, SearchSent)
		coinrest.RegisterAll,

		// Roles createRole handler
		rolerest.RegisterCreateRole,

		// Iris sign transactions handler
		serviceKeys.RegisterSignTx,
		// Iris post transaction handler
		serviceTxs.RegisterPostTx,

		// Iris transfer Tx to byte[]
		serviceByteTx.RegisterByteTx,

		// Nonce query handler
		noncerest.RegisterQueryNonce,

		// Staking query handlers
		stakerest.RegisterQueryCandidate,
		stakerest.RegisterQueryCandidates,
		stakerest.RegisterQueryDelegatorBond,
		stakerest.RegisterQueryDelegatorCandidates,
		// Staking tx builders
		stakerest.RegisterDelegate,
		stakerest.RegisterUnbond,
	}

	for _, routeRegistrar := range routeRegistrars {
		if err := routeRegistrar(router); err != nil {
			log.Fatal(err)
		}
	}

	addr := fmt.Sprintf(":%d", 8080)

	log.Printf("Serving on %q", addr)
	http.ListenAndServe(addr, router)

}