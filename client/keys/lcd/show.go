package keys

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/client/keys/utils"
	"strings"
	"fmt"
)
///////////////////////////
// REST

// get key REST handler
func GetKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	info, err := utils.GetKey(name)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("Key %s not found", name)) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	keyOutput, err := utils.Bech32KeyOutput(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	output, err := json.MarshalIndent(keyOutput, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(output)
}
