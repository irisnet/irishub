package upgrades

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	clientkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

type (
	ClientData struct {
		ChainName      string                 `json:"chain_name"`
		ClientState    map[string]interface{} `json:"client_state"`
		ConsensusState map[string]interface{} `json:"consensus_state"`
		Relayers       []string               `json:"relayers"`
	}

	Client struct {
		ChainName      string
		ClientState    exported.ClientState
		ConsensusState exported.ConsensusState
		Relayers       []string
	}
)

func CreateClient(
	ctx sdk.Context,
	cdc codec.Codec,
	data []byte,
	clientKeeper clientkeeper.Keeper,
) error {
	clients := loadClient(cdc, data)
	for _, client := range clients {
		// init tibc client
		if err := clientKeeper.CreateClient(
			ctx,
			client.ChainName,
			client.ClientState,
			client.ConsensusState,
		); err != nil {
			return err
		}
		// register client relayers
		clientKeeper.RegisterRelayers(ctx, client.ChainName, client.Relayers)
	}
	return nil
}

func loadClient(cdc codec.Codec, data []byte) (clients []Client) {
	var datas []ClientData
	if err := json.Unmarshal(data, &datas); err != nil {
		panic("Unmarshal client.json failed")
	}

	for _, data := range datas {
		clientStateBz, err := json.Marshal(data.ClientState)
		if err != nil {
			panic("Marshal ClientState failed")
		}
		var clientState exported.ClientState
		if err := cdc.UnmarshalInterfaceJSON(clientStateBz, &clientState); err != nil {
			panic("UnmarshalInterfaceJSON ClientState failed")
		}

		consensusStateBz, err := json.Marshal(data.ConsensusState)
		if err != nil {
			panic("Marshal ConsensusState failed")
		}
		var consensusState exported.ConsensusState
		if err := cdc.UnmarshalInterfaceJSON(consensusStateBz, &consensusState); err != nil {
			panic("UnmarshalInterfaceJSON ConsensusState failed")
		}

		clients = append(clients, Client{
			ChainName:      data.ChainName,
			ClientState:    clientState,
			ConsensusState: consensusState,
			Relayers:       data.Relayers,
		})
	}
	return
}
