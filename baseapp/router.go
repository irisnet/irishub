package baseapp

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Router provides handlers for each transaction type.
type Router interface {
	AddRoute(r string, s *sdk.KVStoreKey, h sdk.Handler) (rtr Router)
	Route(path string) (h sdk.Handler)
	RouteTable() (table []string)
}

// map a transaction type to a handler and an initgenesis function
type route struct {
	r string
	s *sdk.KVStoreKey
	h sdk.Handler
}

type router struct {
	routes []route
}

// nolint
// NewRouter - create new router
// TODO either make Function unexported or make return type (router) Exported
func NewRouter() *router {
	return &router{
		routes: make([]route, 0),
	}
}

var isAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

// AddRoute - TODO add description
func (rtr *router) AddRoute(r string, s *sdk.KVStoreKey, h sdk.Handler) Router {
	if !isAlpha(r) {
		panic("route expressions can only contain alphabet characters")
	}
	rtr.routes = append(rtr.routes, route{r,s,h})

	return rtr
}

// Route - TODO add description
// TODO handle expressive matches.
func (rtr *router) Route(path string) (h sdk.Handler) {
	for _, route := range rtr.routes {
		if route.r == path {
			return route.h
		}
	}
	return nil
}

func (rtr *router) RouteTable() (table []string) {
	for _, route := range rtr.routes {
		table = append(table, route.r + "/" + route.s.String())
	}
	return
}
