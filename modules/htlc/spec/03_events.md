<!--
order: 3
-->

# Events

## BeginBlocker

| Type         | Attribute Key | Attribute Value |
| :----------- | :------------ | :-------------- |
| htlc_expired | hash_lock     | {hashLock}      |

## Handlers

### MsgCreateHTLC

| Type        | Attribute Key           | Attribute Value        |
| :---------- | :---------------------- | :--------------------- |
| create_htlc | sender                  | {senderAddress}        |
| create_htlc | receiver                | {receiverAddress}      |
| create_htlc | receiver_on_other_chain | {receiverOnOtherChain} |
| create_htlc | amount                  | {amount}               |
| create_htlc | hash_lock               | {hashLock}             |
| create_htlc | time_lock               | {timeLock}             |
| message     | module                  | htlc                   |
| message     | sender                  | {senderAddress}        |

### MsgClaimHTLC

| Type       | Attribute Key | Attribute Value |
| :--------- | :------------ | :-------------- |
| claim_htlc | sender        | {senderAddress} |
| claim_htlc | hash_lock     | {hashLock}      |
| claim_htlc | secret        | {secret}        |
| message    | module        | htlc            |
| message    | sender        | {senderAddress} |

### MsgRefundHTLC

| Type        | Attribute Key | Attribute Value |
| :---------- | :------------ | :-------------- |
| refund_htlc | sender        | {senderAddress} |
| refund_htlc | hash_lock     | {hashLock}      |
| message     | module        | htlc            |
| message     | sender        | {senderAddress} |
