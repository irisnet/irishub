package keys

import (
	"net/http"

	"github.com/irisnet/irishub/client/keys"
)

// query key list REST handler
func QueryKeysRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kb, err := keys.GetKeyBase()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		infos, err := kb.List()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// an empty list will be JSONized as null, but we want to keep the empty list
		if len(infos) == 0 {
			w.Write([]byte("[]"))
			return
		}
		keysOutput, err := keys.Bech32KeysOutput(infos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		keys.PostProcessResponse(w, cdc, keysOutput, indent)
	}
}
