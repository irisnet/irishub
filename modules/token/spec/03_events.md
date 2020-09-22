<!--
order: 3
-->

# Events

The token module emits the following events:

## Handlers

### MsgIssueToken

| Type        | Attribute Key | Attribute Value |
| ----------- | ------------- | --------------- |
| issue_token | symbol        | {symbol}        |
| message     | module        | token           |
| message     | sender        | {ownerAddress}  |

### MsgEditToken

| Type       | Attribute Key | Attribute Value |
| ---------- | ------------- | --------------- |
| edit_token | symbol        | {symbol}        |
| message    | module        | token           |
| message    | sender        | {ownerAddress}  |

### MsgTransferTokenOwner

| Type                 | Attribute Key | Attribute Value    |
| -------------------- | ------------- | ------------------ |
| transfer_token_owner | validator     | {validatorAddress} |
| delegate             | amount        | {delegationAmount} |
| message              | module        | token              |
| message              | sender        | {ownerAddress}     |

### MsgMintToken

| Type       | Attribute Key | Attribute Value |
| ---------- | ------------- | --------------- |
| mint_token | symbol        | {symbol}        |
| mint_token | amount        | {amount}        |
| message    | module        | token           |
| message    | sender        | {ownerAddress}  |
