package tibc

import (
	_ "embed"
	"encoding/json"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/cosmos/cosmos-sdk/codec"
)

//go:embed clients.json
var clientsByte []byte

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

func LoadClient(cdc codec.Codec) (clients []Client) {
	var datas []ClientData
	if err := json.Unmarshal(clientsByte, &datas); err != nil {
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
