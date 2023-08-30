# NFT

## Introduction

`NFT` provides the ability to digitize assets. Through this module, each off-chain asset will be modeled as a unique on-chain nft.

nft on the chain are identified by `ID`. With the help of the secure and non-tamperable features of the blockchain, the ownership of nft will be clarified. The transaction process of nft among members will also be publicly recorded to facilitate traceability and dispute settlement.

nft metadata (`metadata`) can be stored directly on the chain, or the `URI` of its storage source outside the chain can be stored on the chain. nft metadata is organized according to a specific [JSON Schema](https://JSON-Schema.org/). [Here](https://github.com/irisnet/irishub/blob/master/docs/features/nft-metadata.json) is an example of metadata JSON Schema.

nft need to be issued before creation to declare their abstract properties:

-_Denom_: the globally unique nft classification name

-_Denom ID_: the globally unique nft classification identifier of Denom

-_Symbol_: the symbol of the token

-_Mint-restricted_: This field indicates whether there are restrictions on the issuance of NFTs under this classification, true means that only Denom owners can issue NFTs under this classification, false means anyone can

-_Update-restricted_: This field indicates whether there are restrictions on updating NFTs under this classification, true means that no one under this classification can update the NFT, false means that only the owner of this NFT can update

-_Metadata Specification_: The JSON Schema that nft metadata should follow

Each specific nft is described by the following elements:

-_Denom_: the classification of the nft

-_ID_: The identifier of the nft, which is unique in this nft denom; this ID is generated off-chain

-_Metadata_: The structure containing the specific data of the nft

-_Metadata URI_: When metadata is stored off-chain, this URI indicates its storage location

## Features

### issued

Specify the nft Denom (nft classification) and metadata JSON Schema to issue nft.

`CLI`

```bash
iris tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --schema=<schema-content or path/to/schema.json> --symbol=<denom-symbol> --mint-restricted=<mint-restricted>  --update-restricted=<update-restricted> --chain-id=<chain-id> --fees=<fee>
```

### transfer denom

The owner of the NFT classification can transfer the ownership of the NFT classification to others.

`CLI`

```bash
iris tx nft transfer-denom <recipient> <denom-id>
```

### Additional issuance

After the nft is issued, additional issuance (create) of specific nft of this type can be made. The denom ID, token ID, recipient address and URI must be specified.

`CLI`

```bash
iris tx nft mint <denom-id> <token-id> --uri=<uri> --recipient=<recipient> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### Edit

The metadata of the specified nft can be updated.

`CLI`

```bash
iris tx nft edit <denom-id> <token-id> --uri=<uri> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### Transfer

Transfer designated nft.

`CLI`

```bash
iris tx nft transfer <recipient-address> <denom-id> <token-id>
```

### Destroy

You can destroy the created nft.

`CLI`

```bash
iris tx nft burn <denom-id> <token-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### Query the specified nft denom

Query nft denom information based on Denom ID.

`CLI`

```bash
iris q nft denom <denom-id>
```

### Query all nft denom information

Query all issued nft denom information.

`CLI`

```bash
iris q nft denoms
```

### Query the total amount of nft in a specified denom

Query the total amount of nft according to Denom ID; accept the optional owner parameter.

`CLI`

```bash
iris q nft supply <denom-id> --owner=<owner>
```

### Query all nft of the specified account

Query all nft owned by an account; you can specify the Denom ID parameter.

`CLI`

```bash
iris q nft owner <address> --denom-id=<denom-id>
```

### Query all nft of a specified denom

Query all nft according to Denom ID.

`CLI`

```bash
iris q nft collection <denom-id>
```

### Query specified nft

Query specific nft based on Denom ID and Token ID.

`CLI`

```bash
iris q nft token <denom-id> <token-id>
```
