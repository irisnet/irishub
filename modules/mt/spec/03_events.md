# Events

The mt module emits the following events:

## Handlers

### MsgIssueDenom

| Type        | Attribute Key | Attribute Value  |
|:------------|:--------------|:-----------------|
| issue_denom | denom_id      | {mtDenomID}     |
| issue_denom | denom_name    | {mtDenomName}   |
| issue_denom | creator       | {creatorAddress} |
| message     | module        | mt              |
| message     | sender        | {senderAddress}  |

### MsgTransferMT

| Type         | Attribute Key | Attribute Value    |
|:-------------|:--------------|:-------------------|
| transfer_mt | token_id      | {tokenID}          |
| transfer_mt | denom_id      | {mtDenomID}       |
| transfer_mt | sender        | {senderAddress}    |
| transfer_mt | recipient     | {recipientAddress} |
| message      | module        | mt                |
| message      | sender        | {senderAddress}    |

### MsgEditMT

| Type     | Attribute Key | Attribute Value |
|:---------|:--------------|:----------------|
| edit_mt | token_id      | {tokenID}       |
| edit_mt | denom_id      | {mtDenomID}    |
| edit_mt | token_uri     | {tokenURI}      |
| edit_mt | owner         | {ownerAddress}  |
| message  | module        | mt             |
| message  | sender        | {senderAddress} |

### MsgMintMT

| Type     | Attribute Key | Attribute Value    |
|:---------|:--------------|:-------------------|
| mint_mt | token_id      | {tokenID}          |
| mint_mt | denom_id      | {mtDenomID}       |
| mint_mt | token_uri     | {tokenURI}         |
| mint_mt | recipient     | {recipientAddress} |
| message  | module        | mt                |
| message  | sender        | {senderAddress}    |

### MsgBurnMTs

| Type     | Attribute Key | Attribute Value |
|:---------|:--------------|:----------------|
| burn_mt | denom_id      | {mtDenomID}    |
| burn_mt | token_id      | {tokenID}       |
| burn_mt | owner         | {ownerAddress}  |
| message  | module        | mt             |
| message  | sender        | {senderAddress} |

### MsgTransferDenom

| Type           | Attribute Key | Attribute Value    |
|:---------------|:--------------|:-------------------|
| transfer_denom | denom_id      | {mtDenomID}       |
| transfer_denom | sender        | {senderAddress}    |
| transfer_denom | recipient     | {recipientAddress} |
| message        | module        | mt                |
| message        | sender        | {senderAddress}    |