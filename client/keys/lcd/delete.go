package keys

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	"net/http"
	"strings"
	"github.com/irisnet/irishub/client/utils"
)

// delete key request REST body
type DeleteKeyBody struct {
	Password string `json:"password"`
}

// delete key REST handler
func DeleteKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var m DeleteKeyBody

	err := utils.ReadPostBody(w, r, cdc, &m)
	if err != nil {
		return
	}

	kb, err := keys.GetKeyBase()
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
