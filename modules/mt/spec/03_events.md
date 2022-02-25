# Events

The mt module emits the following events:

## Handlers

### MsgIssueDenom

| Type        | Attribute Key | Attribute Value  |
|:------------|:--------------|:-----------------|
| issue_denom | denom_id      | {mtDenomID}      |
| issue_denom | denom_name    | {mtDenomName}    |
| issue_denom | owner         | {senderAddress}  |
| message     | module        | mt               |
| message     | sender        | {senderAddress}  |

### MsgTransferDenom

| Type           | Attribute Key | Attribute Value    |
|:---------------|:--------------|:-------------------|
| transfer_denom | denom_id      | {mtDenomID}        |
| transfer_denom | recipient     | {recipientAddress} |
| message        | module        | mt                 |
| message        | sender        | {senderAddress}    |

### MsgMintMT

| Type     | Attribute Key | Attribute Value    |
|:---------|:--------------|:-------------------|
| mint_mt  | mt_id         | {mtID}             |
| mint_mt  | denom_id      | {mtDenomID}        |
| mint_mt  | amount        | {mintAmount}       |
| mint_mt  | supply        | {mtSupply}         |
| mint_mt  | recipient     | {recipientAddress} |
| message  | module        | mt                 |
| message  | sender        | {senderAddress}    |

### MsgEditMT

| Type     | Attribute Key | Attribute Value |
|:---------|:--------------|:----------------|
| edit_mt  | mt_id         | {mtID}          |
| edit_mt  | denom_id      | {mtDenomID}     |
| message  | module        | mt              |
| message  | sender        | {senderAddress} |

### MsgTransferMT

| Type         | Attribute Key | Attribute Value    |
|:-------------|:--------------|:-------------------|
| transfer_mt  | mt_id         | {mtID}             |
| transfer_mt  | denom_id      | {mtDenomID}        |
| transfer_mt  | amount        | {transferAmount}   |
| transfer_mt  | recipient     | {recipientAddress} |
| message      | module        | mt                 |
| message      | sender        | {senderAddress}    |

### MsgBurnMTs

| Type     | Attribute Key | Attribute Value |
|:---------|:--------------|:----------------|
| burn_mt  | mt_id         | {mtID}          |
| burn_mt  | denom_id      | {mtDenomID}     |
| burn_mt  | amount        | {burnAmount}    |
| message  | module        | mt              |
| message  | sender        | {senderAddress} |
