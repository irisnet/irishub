<!--
order: 3
-->

# Events

The service module emits the following events:

## EndBlocker

| Type                       | Attribute Key         | Attribute Value       |
|:---------------------------|:----------------------|:----------------------|
| new_batch                  | request_context_id    | {requestContextID}    |
| new_batch                  | request_context_state | {requestContextState} |
| new_batch_request_provider | service_name          | {serviceName}         |
| new_batch_request_provider | provider              | {providerAddress}     |
| new_batch_request_provider | requests              | {requests}            |
| new_batch_request          | service_name          | {serviceName}         |
| new_batch_request          | request_context_id    | {requestContextID}    |
| new_batch_request          | requests              | {requests}            |
| complete_batch             | request_context_id    | {requestContextID}    |
| complete_batch             | request_context_state | {requestContextState} |
| service_slash              | slashed_coins         | {slashedCoins}        |
| service_slash              | provider              | {providerAddress}     |
| service_slash              | requests              | {requests}            |
| pause_context              | request_context_id    | {requestContextID}    |
| complete_context           | request_context_id    | {requestContextID}    |
| no_exchange_rate           | price_denom           | {priceDenom}          |
| no_exchange_rate           | request_context_id    | {requestContextID}    |
| no_exchange_rate           | service_name          | {serviceName}         |
| no_exchange_rate           | consumer              | {consumerAddress}     |

## Handlers

### MsgDefineService

| Type              | Attribute Key | Attribute Value |
|:------------------|:--------------|:----------------|
| create_definition | service_name  | {serviceName}   |
| create_definition | author        | {authorAddress} |
| message           | module        | service         |
| message           | sender        | {senderAddress} |

### MsgBindService

| Type           | Attribute Key | Attribute Value   |
|:---------------|:--------------|:------------------|
| create_binding | service_name  | {serviceName}     |
| create_binding | provider      | {providerAddress} |
| create_binding | owner         | {ownerAddress}    |
| message        | module        | service           |
| message        | sender        | {senderAddress}   |

### MsgUpdateServiceBinding

| Type           | Attribute Key | Attribute Value   |
|:---------------|:--------------|:------------------|
| update_binding | service_name  | {serviceName}     |
| update_binding | provider      | {providerAddress} |
| update_binding | owner         | {ownerAddress}    |
| message        | module        | service           |
| message        | sender        | {senderAddress}   |

### MsgDisableServiceBinding

| Type            | Attribute Key | Attribute Value   |
|:----------------|:--------------|:------------------|
| disable_binding | service_name  | {serviceName}     |
| disable_binding | provider      | {providerAddress} |
| disable_binding | owner         | {ownerAddress}    |
| message         | module        | service           |
| message         | sender        | {senderAddress}   |

### MsgEnableServiceBinding

| Type           | Attribute Key | Attribute Value   |
|:---------------|:--------------|:------------------|
| enable_binding | service_name  | {serviceName}     |
| enable_binding | provider      | {providerAddress} |
| enable_binding | owner         | {ownerAddress}    |
| message        | module        | service           |
| message        | sender        | {senderAddress}   |

### MsgRefundServiceDeposit

| Type           | Attribute Key | Attribute Value   |
|:---------------|:--------------|:------------------|
| refund_deposit | service_name  | {serviceName}     |
| refund_deposit | provider      | {providerAddress} |
| refund_deposit | owner         | {ownerAddress}    |
| message        | module        | service           |
| message        | sender        | {senderAddress}   |

### MsgSetWithdrawAddress

| Type                 | Attribute Key    | Attribute Value   |
|:---------------------|:-----------------|:------------------|
| set_withdraw_address | withdraw_address | {withdrawAddress} |
| set_withdraw_address | owner            | {ownerAddress}    |
| message              | module           | service           |
| message              | sender           | {senderAddress}   |

### MsgCallService

| Type           | Attribute Key      | Attribute Value    |
|:---------------|:-------------------|:-------------------|
| create_context | request_context_id | {requestContextID} |
| create_context | service_name       | {serviceName}      |
| create_context | consumer           | {consumerAddress}  |
| message        | module             | service            |
| message        | sender             | {senderAddress}    |

### MsgRespondService

| Type            | Attribute Key      | Attribute Value    |
|:----------------|:-------------------|:-------------------|
| respond_service | request_context_id | {requestContextID} |
| respond_service | request_id         | {requestID}        |
| respond_service | service_name       | {serviceName}      |
| respond_service | provider           | {providerAddress}  |
| respond_service | consumer           | {consumerAddress}  |
| message         | module             | service            |
| message         | sender             | {senderAddress}    |

### MsgUpdateRequestContext

| Type           | Attribute Key      | Attribute Value    |
|:---------------|:-------------------|:-------------------|
| update_context | request_context_id | {requestContextID} |
| update_context | consumer           | {consumerAddress}  |
| message        | module             | service            |
| message        | sender             | {senderAddress}    |

### MsgPauseRequestContext

| Type          | Attribute Key      | Attribute Value    |
|:--------------|:-------------------|:-------------------|
| pause_context | request_context_id | {requestContextID} |
| pause_context | consumer           | {consumerAddress}  |
| message       | module             | service            |
| message       | sender             | {senderAddress}    |

### MsgStartRequestContext

| Type          | Attribute Key      | Attribute Value    |
|:--------------|:-------------------|:-------------------|
| start_context | request_context_id | {requestContextID} |
| start_context | consumer           | {consumerAddress}  |
| message       | module             | service            |
| message       | sender             | {senderAddress}    |

### MsgKillRequestContext

| Type         | Attribute Key      | Attribute Value    |
|:-------------|:-------------------|:-------------------|
| kill_context | request_context_id | {requestContextID} |
| kill_context | consumer           | {consumerAddress}  |
| message      | module             | service            |
| message      | sender             | {senderAddress}    |

### MsgWithdrawEarnedFees

| Type                 | Attribute Key | Attribute Value   |
|:---------------------|:--------------|:------------------|
| withdraw_earned_fees | provider      | {providerAddress} |
| withdraw_earned_fees | owner         | {ownerAddress}    |
| message              | module        | service           |
| message              | sender        | {senderAddress}   |

