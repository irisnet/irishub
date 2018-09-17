package keys

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/client/keys/utils"
	"strings"
)

////////////////////////
// REST

// delete key request REST body
type DeleteKeyBody struct {
	Password string `json:"password"`
}

// delete key REST handler
func DeleteKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m DeleteKeyBody

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

	err = kb.Delete(name, m.Password)
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
