// nolint
package types

// guardian module event types
const (
	EventTypeAddProfiler    = "add_profiler"
	EventTypeAddTrustee     = "add_trustee"
	EventTypeDeleteProfiler = "delete_profiler"
	EventTypeDeleteTrustee  = "delete_trustee"

	AttributeKeyProfilerAddress = "address"
	AttributeKeyTrusteeAddress  = "address"
	AttributeKeyAddedBy         = "added_by"
	AttributeKeyDeletedBy       = "deleted_by"

	AttributeValueCategory = ModuleName
)
