package simapp

import "cosmossdk.io/depinject"

// DepinjectOptions are passed to the app on creation
type DepinjectOptions struct {
	Config    depinject.Config
	Providers []interface{}
	Consumers []interface{}
}

