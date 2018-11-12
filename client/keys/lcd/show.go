package keys

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	keycli "github.com/irisnet/irishub/client/keys/cli"
)

///////////////////////////
// REST

// get key REST handler
func GetKeyRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		bechPrefix := r.URL.Query().Get(keycli.FlagBechPrefix)

		if bechPrefix == "" {
			bechPrefix = "acc"
		}

		bechKeyOut, err := keys.GetBechKeyOut(bechPrefix)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		info, err := keys.GetKeyInfo(name)
		if keyerror.IsErrKeyNotFound(err) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(err.Error()))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		keyOutput, err := bechKeyOut(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		keys.PostProcessResponse(w, cdc, keyOutput, indent)
	}
}
