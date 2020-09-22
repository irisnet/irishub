<!--
order: 3
-->

# Events

The service module emits the following events:

## EndBlocker

| Type                       | Attribute Key      | Attribute Value    |
| -------------------------- | ------------------ | ------------------ |
| new_batch                  | request_context_id | {requestContextID} |
| new_batch_request_provider | requests           | {requests}         |
| complete_batch             | batch_state        | {batchState}       |
| pause_context              | request_context_id | {requestContextID} |
| service_slash              | slashed_coins      | {slashedCoins}     |

## Handlers

### MsgDefineService

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgBindService

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgUpdateServiceBinding

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgDisableServiceBinding

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgEnableServiceBinding

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgRefundServiceDeposit

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgSetWithdrawAddress

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgCallService

| Type           | Attribute Key      | Attribute Value    |
| -------------- | ------------------ | ------------------ |
| create_context | request_context_id | {requestContextID} |
| message        | module             | service            |
| message        | sender             | {senderAddress}    |

### MsgRespondService

| Type            | Attribute Key | Attribute Value |
| --------------- | ------------- | --------------- |
| respond_service | request_id    | {requestID}     |
| message         | module        | service         |
| message         | sender        | {senderAddress} |

### MsgUpdateRequestContext

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgPauseRequestContext

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgStartRequestContext

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgKillRequestContext

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |

### MsgWithdrawEarnedFees

| Type    | Attribute Key | Attribute Value |
| ------- | ------------- | --------------- |
| message | module        | service         |
| message | sender        | {senderAddress} |
