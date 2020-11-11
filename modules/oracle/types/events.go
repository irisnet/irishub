package types

// guardian module event types
const (
	EventTypeCreateFeed   = "create_feed"
	EventTypeStartFeed    = "start_feed"
	EventTypePauseFeed    = "pause_feed"
	EventTypeEditFeed     = "edit_feed"
	EventTypeSetFeedValue = "set_feed_value"

	AttributeValueCategory = ModuleName

	AttributeKeyFeedName    = "feed_name"
	AttributeKeyServiceName = "service_name"
	AttributeKeyFeedValue   = "feed_value"
	AttributeKeyCreator     = "creator"
)
