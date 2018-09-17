package keys

import (
	"github.com/gorilla/mux"
)

// resgister REST routes
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/keys", QueryKeysRequestHandler).Methods("GET")
	r.HandleFunc("/keys", AddNewKeyRequestHandler).Methods("POST")
	r.HandleFunc("/keys/{name}/sign", GetSignRequestHandler).Methods("POST")
	r.HandleFunc("/keys/{name}", GetKeyRequestHandler).Methods("GET")
	r.HandleFunc("/keys/{name}", UpdateKeyRequestHandler).Methods("PUT")
	r.HandleFunc("/keys/{name}", DeleteKeyRequestHandler).Methods("DELETE")
}
