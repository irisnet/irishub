package keys

import (
	"io/ioutil"
	"net/http"

	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/client/utils"
)

// NewKeyBody - the request body for create or recover new keys
type NewKeyBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

// AddNewKeyRequestHandler performs create or recover new keys operation
func AddNewKeyRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var kb cryptokeys.Keybase
		var m NewKeyBody

		kb, err := keys.GetKeyBase()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}

		if m.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = keys.ErrMissingName()
			w.Write([]byte(err.Error()))
			return
		}
		if m.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = keys.ErrMissingPassword()
			w.Write([]byte(err.Error()))
			return
		}

		// check if already exists
		infos, err := kb.List()
		for _, i := range infos {
			if i.GetName() == m.Name {
				w.WriteHeader(http.StatusConflict)
				err = keys.ErrKeyNameConflict(m.Name)
				w.Write([]byte(err.Error()))
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

		keys.PostProcessResponse(w, cdc, keyOutput, indent)
	}
}

// function to just a new seed to display in the UI before actually persisting it in the keybase
func getSeed(algo cryptokeys.SigningAlgo) string {
	kb := keys.MockKeyBase()
	pass := "throwing-this-key-away"
	name := "inmemorykey"
	_, seed, _ := kb.CreateMnemonic(name, cryptokeys.English, pass, algo)
	return seed
}

// Seed REST request handler
func SeedRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	algoType := vars["type"]
	// algo type defaults to secp256k1
	if algoType == "" {
		algoType = "secp256k1"
	}
	algo := cryptokeys.SigningAlgo(algoType)

	seed := getSeed(algo)

	w.Write([]byte(seed))
}

// RecoverKeyBody is recover key request REST body
type RecoverKeyBody struct {
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

// RecoverRequestHandler performs key recover request
func RecoverRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		var m RecoverKeyBody
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

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = keys.ErrMissingName()
			w.Write([]byte(err.Error()))
			return
		}
		if m.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = keys.ErrMissingPassword()
			w.Write([]byte(err.Error()))
			return
		}
		if m.Seed == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = keys.ErrMissingSeed()
			w.Write([]byte(err.Error()))
			return
		}

		kb, err := keys.GetKeyBaseWithWritePerm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// check if already exists
		infos, err := kb.List()
		for _, info := range infos {
			if info.GetName() == name {
				w.WriteHeader(http.StatusConflict)
				err = keys.ErrKeyNameConflict(name)
				w.Write([]byte(err.Error()))
				return
			}
		}

		info, err := kb.CreateKey(name, m.Seed, m.Password)
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

		keys.PostProcessResponse(w, cdc, keyOutput, indent)
	}
}
