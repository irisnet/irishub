package exported

import (
	"github.com/irisnet/irismod/modules/service/types"
)

type (
	// RequestContext defines a context which holds request-related data
	RequestContext = types.RequestContext

	// RequestContextState represents the request context state
	RequestContextState = types.RequestContextState

	// ResponseCallback defines the response callback interface
	ResponseCallback = types.ResponseCallback

	// StateCallback defines the state callback interface
	StateCallback = types.StateCallback

	// ModuleService defines the moduler service interface
	ModuleService = types.ModuleService
)

const (
	// RUNNING indicates the request context is running
	RUNNING = types.RUNNING

	// PAUSED indicates the request context is paused
	PAUSED = types.PAUSED

	// COMPLETED indicates the request context is completed
	COMPLETED = types.COMPLETED
)

var (
	RequestContextStateFromString = types.RequestContextStateFromString
)
