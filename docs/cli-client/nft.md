# NFT

`NFT` provides the ability to digitize assets. Through this module, each off-chain asset will be modeled as a unique on-chain nft.

## Available Commands

| Name                                          | Description                                                                                         |
| --------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [issue](#iris-tx-nft-issue)                   | Specify the nft Denom (nft classification) and metadata JSON Schema to issue nft.                   |
| [transfer-denom](#iris-tx-nft-transfer-denom) | The owner of the NFT classification can transfer the ownership of the NFT classification to others. |
| [mint](#iris-tx-nft-mint)                     | Additional issuance (create) of specific nft of this type can be made.                              |
| [edit](#iris-tx-nft-edit)                     | The metadata of the specified nft can be updated.                                                   |
| [transfer](#iris-tx-nft-transfer)             | Transfer designated nft.                                                                            |
| [burn](#iris-tx-nft-burn)                     | Destroy the created nft.                                                                            |
| [supply](#iris-query-nft-supply)              | Query the total amount of nft according to Denom; accept the optional owner parameter.              |
| [owner](#iris-query-nft-owner)                | Query all nft owned by an account; you can specify the Denom parameter.                             |
| [collection](#iris-query-nft-collection)      | Query all nft according to Denom.                                                                   |
| [denom](#iris-query-nft-denom)                | Query nft denom information based on Denom.                                                         |
| [denoms](#iris-query-nft-denoms)              | Query the total amount of nft according to Denom; accept the optional owner parameter.              |
| [token](#iris-query-nft-token)                | Query specific nft based on Denom and ID.                                                           |

## iris tx nft issue

Specify the nft Denom (nft classification) and metadata JSON Schema to issue nft.

```bash
iris tx nft issue [denom-id] [flags]
```

**Flags:**

| Name, shorthand     | Required | Default                                                                                                                                                                                                                     | Description |
| ------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| --name              |          | The name of the denom                                                                                                                                                                                                       |             |
| --uri               |          | The uri of the denom                                                                                                                                                                                                        |             |
| --data              |          | Off-chain metadata for supplementation (JSON object)                                                                                                                                                                        |             |
| --schema            |          | Denom data structure definition                                                                                                                                                                                             |             |
| --symbol            |          | The symbol of the denom                                                                                                                                                                                                     |             |
| --mint-restricted   |          | This field indicates whether there are restrictions on the issuance of NFTs under this classification, true means that only Denom owners can issue NFTs under this classification, false means anyone can                   |             |
| --update-restricted |          | This field indicates whether there are restrictions on updating NFTs under this classification, true means that no one under this classification can update the NFT, false means that only the owner of this NFT can update |             |

## iris tx nft transfer-denom

The owner of the NFT classification can transfer the ownership of the NFT classification to others.

```bash
iris tx nft transfer-denom [recipient] [denom-id]
```

## iris tx nft mint

Additional issuance (create) of specific nft of this type can be made.  

```bash
iris tx nft mint [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --uri           |          | URI of off-chain token data |             |
| --recipient     |          | Receiver of the nft         |             |
| --name          |          | The name of nft             |             |

## iris tx nft edit

The metadata of the specified nft can be updated.

```bash
iris tx nft edit [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --uri           |          | URI of off-chain token data |             |
| --name          |          | The name of nft             |             |

## iris tx nft transfer

Transfer designated nft.

```bash
iris tx nft transfer [recipient] [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --uri           |          | URI of off-chain token data |             |
| --name          |          | The name of nft             |             |

## iris tx nft burn

Destroy the created nft.

```bash
iris tx nft burn [denomID] [tokenID] [flags]
```

## iris query nft

Query nft

### iris query nft supply

```bash
iris query nft supply [denomID]
iris query nft supply [denomID] --owner=<owner address>
```

### iris query nft owner

```bash
iris query nft owner [owner address]
iris query nft owner [owner address] --denom=<denomID>
```

### iris query nft collection

```bash
iris query nft collection [denomID]
```

### iris query nft denom

```bash
iris query nft denom [denomID]
```

### iris query nft denoms

```bash
iris query nft denoms
```

### iris query nft token

```bash
iris query nft token [denomID] [tokenID]
```
