package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
)

// NewKeyBody - the request body for create or recover new keys
type NewKeyBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

// AddNewKeyRequestHandler performs create or recover new keys operation
func AddNewKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	var kb cryptokeys.Keybase
	var m NewKeyBody

	kb, err := keys.GetKeyBase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = cdc.UnmarshalJSON(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if m.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You have to specify a name for the locally stored account."))
		return
	}
	if m.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You have to specify a password for the locally stored account."))
		return
	}

	// check if already exists
	infos, err := kb.List()
	for _, i := range infos {
		if i.GetName() == m.Name {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("Account with name %s already exists.", m.Name)))
			return
		}
	}

	// create account
	seed := m.Seed
	if seed == "" {
		seed = getSeed(cryptokeys.Secp256k1)
	}
	info, err := kb.CreateKey(m.Name, seed, m.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	keyOutput, err := keys.Bech32KeyOutput(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	keyOutput.Seed = seed

	bz, err := json.Marshal(keyOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bz)
}

// function to just a new seed to display in the UI before actually persisting it in the keybase
func getSeed(algo cryptokeys.SigningAlgo) string {
	kb := keys.MockKeyBase()
	pass := "throwing-this-key-away"
	name := "inmemorykey"
	_, seed, _ := kb.CreateMnemonic(name, cryptokeys.English, pass, algo)
	return seed
}
