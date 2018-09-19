package keys

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	"strings"
	"github.com/irisnet/irishub/client/utils"
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
	var m UpdateKeyBody

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
