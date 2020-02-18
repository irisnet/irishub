package exported

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
)

type (
	// RequestContext defines a context which holds request-related data
	RequestContext types.RequestContext

	// RequestContextState represents the request context state
	RequestContextState types.RequestContextState

	// ResponseCallback defines the response callback interface
	ResponseCallback types.ResponseCallback
)

const (
	// RUNNING indicates the request context is running
	RUNNING = types.RUNNING

	// PAUSED indicates the request context is paused
	PAUSED = types.PAUSED
)
