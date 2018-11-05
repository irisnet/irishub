package lcd

import (
	"fmt"
	"net/http"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/version"
)

// cli version REST handler endpoint
func CLIVersionRequestHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()
	w.Write([]byte(v))
}

// connected node version REST handler endpoint
func NodeVersionRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := cliCtx.Query("/app/version",nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query version. Error: %s", err.Error())))
			return
		}

		w.Write(version)
	}
}