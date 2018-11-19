package keys

import (
	"github.com/gorilla/mux"
)

// resgister REST routes
func RegisterRoutes(r *mux.Router, indent bool) {
	r.HandleFunc("/keys", QueryKeysRequestHandler(indent)).Methods("GET")
	r.HandleFunc("/keys", AddNewKeyRequestHandler(indent)).Methods("POST")
	r.HandleFunc("/keys/seed", SeedRequestHandler(indent)).Methods("GET")
	r.HandleFunc("/keys/{name}/recover", RecoverRequestHandler(indent)).Methods("POST")
	r.HandleFunc("/keys/{name}", GetKeyRequestHandler(indent)).Methods("GET")
	r.HandleFunc("/keys/{name}", UpdateKeyRequestHandler).Methods("PUT")
	r.HandleFunc("/keys/{name}", DeleteKeyRequestHandler).Methods("DELETE")
}
