<!--
order: 3
-->

# Events

## Handlers

### MsgCreateFeed

| Type        | Attribute Key | Attribute Value  |
| :---------- | :------------ | :--------------- |
| create_feed | feed_name     | {feedName}       |
| create_feed | service_name  | {serviceName}    |
| create_feed | creator       | {creatorAddress} |
| message     | module        | oracle           |
| message     | sender        | {senderAddress}  |

### MsgStartFeed

| Type       | Attribute Key | Attribute Value  |
| :--------- | :------------ | :--------------- |
| start_feed | feed_name     | {feedName}       |
| start_feed | creator       | {creatorAddress} |
| message    | module        | oracle           |
| message    | sender        | {senderAddress}  |

### MsgPauseFeed

| Type       | Attribute Key | Attribute Value  |
| :--------- | :------------ | :--------------- |
| pause_feed | feed_name     | {feedName}       |
| pause_feed | creator       | {creatorAddress} |
| message    | module        | oracle           |
| message    | sender        | {senderAddress}  |

### MsgEditFeed

| Type      | Attribute Key | Attribute Value  |
| :-------- | :------------ | :--------------- |
| edit_feed | feed_name     | {feedName}       |
| edit_feed | creator       | {creatorAddress} |
| message   | module        | oracle           |
| message   | sender        | {senderAddress}  |

## Callbacks

| Type           | Attribute Key | Attribute Value |
| :------------- | :------------ | :-------------- |
| set_feed_value | feed_name     | {feedName}      |
| set_feed_value | feed_value    | {feedValue}     |
