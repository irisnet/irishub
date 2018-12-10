package protocol

import (
	"regexp"

	sdk "github.com/irisnet/irishub/types"
	"strings"
)

// Router provides handlers for each transaction type.
type Router interface {
	AddRoute(r string, h sdk.Handler) (rtr Router)
	Route(path string) (h sdk.Handler)
}

// map a transaction type to a handler and an initgenesis function
type route struct {
	r string
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

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// AddRoute - TODO add description
func (rtr *router) AddRoute(r string, h sdk.Handler) Router {
	rstrs := strings.Split(r, "-")

	if !isAlphaNumeric(rstrs[0]) {
		panic("route expressions can only contain alphabet characters")
	}
	rtr.routes = append(rtr.routes, route{r, h})

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
