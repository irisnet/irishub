# Events

The nft module emits the following events:

## Handlers

### MsgIssueDenom

| Type         | Attribute Key | Attribute Value |
| ------------ | ------------- | --------------- |
| transfer_nft | denom         | {nftDenom}      |
| message      | module        | nft             |
| message      | action        | issue_nft       |
| message      | sender        | {senderAddress} |

### MsgTransferNFT

| Type         | Attribute Key | Attribute Value    |
| ------------ | ------------- | ------------------ |
| transfer_nft | denom         | {nftDenom}         |
| transfer_nft | token-id      | {tokenID}          |
| transfer_nft | recipient     | {recipientAddress} |
| message      | module        | nft                |
| message      | action        | transfer_nft       |
| message      | sender        | {senderAddress}    |

### MsgEditNFT

| Type     | Attribute Key | Attribute Value |
| -------- | ------------- | --------------- |
| edit_nft | denom         | {nftDenom}      |
| edit_nft | token-id      | {tokenID}       |
| message  | module        | nft             |
| message  | action        | edit_nft        |
| message  | sender        | {senderAddress} |
| message  | token-uri     | {tokenURI}      |

### MsgMintNFT

| Type     | Attribute Key | Attribute Value |
| -------- | ------------- | --------------- |
| mint_nft | denom         | {nftDenom}      |
| mint_nft | token-id      | {tokenID}       |
| message  | module        | nft             |
| message  | action        | mint_nft        |
| message  | sender        | {senderAddress} |
| message  | token-uri     | {tokenURI}      |

### MsgBurnNFTs

| Type     | Attribute Key | Attribute Value |
| -------- | ------------- | --------------- |
| burn_nft | denom         | {nftDenom}      |
| burn_nft | token-id      | {tokenID}       |
| message  | module        | nft             |
| message  | action        | burn_nft        |
| message  | sender        | {senderAddress} |
