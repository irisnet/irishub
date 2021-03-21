<!--
order: 3
-->

# Events

## BeginBlocker

| Type        | Attribute Key | Attribute Value |
| :---------- | :------------ | :-------------- |
| refund_htlc | id            | {htlcID}        |

## Handlers

### MsgCreateHTLC

| Type        | Attribute Key           | Attribute Value        |
| :---------- | :---------------------- | :--------------------- |
| create_htlc | id                      | {htlcID}               |
| create_htlc | sender                  | {senderAddress}        |
| create_htlc | receiver                | {receiverAddress}      |
| create_htlc | receiver_on_other_chain | {receiverOnOtherChain} |
| create_htlc | sender_on_other_chain   | {senderOnOtherChain}   |
| create_htlc | transfer                | `true`/`false`         |
| message     | module                  | htlc                   |
| message     | sender                  | {senderAddress}        |

### MsgClaimHTLC

| Type       | Attribute Key | Attribute Value |
| :--------- | :------------ | :-------------- |
| claim_htlc | id            | {htlcID}        |
| claim_htlc | hash_lock     | {hashLock}      |
| claim_htlc | sender        | {senderAddress} |
| claim_htlc | secret        | {secret}        |
| claim_htlc | transfer      | `true`/`false`  |
| claim_htlc | direction     | {direction}     |
| message    | module        | htlc            |
| message    | sender        | {senderAddress} |
