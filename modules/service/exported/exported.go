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

	// ModuleService defines the module service interface
	ModuleService = types.ModuleService
)

const (
	// RUNNING indicates the request context is running
	RUNNING = types.RUNNING

	// PAUSED indicates the request context is paused
	PAUSED = types.PAUSED

	// COMPLETED indicates the request context is completed
	COMPLETED = types.COMPLETED

	// RegisterModuleName exports types.RegisterModuleName
	RegisterModuleName = types.RegisterModuleName

	// OraclePriceServiceName exports types.OraclePriceServiceName
	OraclePriceServiceName = types.OraclePriceServiceName

	// PATH_BODY exports types.PATH_BODY
	PATH_BODY = types.PATH_BODY
)

var (
	// RequestContextStateFromString exports types.RequestContextStateFromString
	RequestContextStateFromString = types.RequestContextStateFromString
	// ValidateServiceName exports types.ValidateServiceName
	ValidateServiceName = types.ValidateServiceName
	// OraclePriceServiceProvider exports types.OraclePriceServiceProvider
	OraclePriceServiceProvider = types.OraclePriceServiceProvider
)
