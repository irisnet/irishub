package keys

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	"net/http"
	"github.com/irisnet/irishub/client/utils"
)

type keySignBody struct {
	Tx       []byte `json:"tx"`
	Password string `json:"password"`
}

// GetSignRequestHandler is the handler of creating seed in swagger rest server
func GetSignRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var m keySignBody
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

	sig, _, err := kb.Sign(name, m.Password, m.Tx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encoded := base64.StdEncoding.EncodeToString(sig)

	w.Write([]byte(encoded))
}
