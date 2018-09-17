package keys

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys/utils"
	"strings"
)

// update key request REST body
type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

// update key REST handler
func UpdateKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m UpdateKeyBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	kb, err = utils.GetKeyBase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	getNewpass := func() (string, error) { return m.NewPassword, nil }

	err = kb.Update(name, m.OldPassword, getNewpass)
	if err != nil {
		if strings.Contains(err.Error(), "Ciphertext decryption failed") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
