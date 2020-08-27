# NFT

## Introduction

`NFT` provides the ability to digitize assets. Through this module, each off-chain asset will be modeled as a unique on-chain asset.

Assets on the chain are identified by `ID`. With the help of the secure and non-tamperable features of the blockchain, the ownership of assets will be clarified. The transaction process of assets among members will also be publicly recorded to facilitate traceability and dispute settlement.

Asset metadata (`metadata`) can be stored directly on the chain, or the `URI` of its storage source outside the chain can be stored on the chain. Asset metadata is organized according to a specific [JSON Schema](https://JSON-Schema.org/). [Here](./nft-metadata.json) is an example of metadata JSON Schema.

Assets need to be issued before creation to declare their abstract properties:

-_Denom_: the globally unique asset class identifier

-_Metadata Specification_: The JSON Schema that asset metadata should follow

Each specific asset is described by the following elements:

-_Denom_: the category of the asset

-_ID_: The identifier of the asset, which is unique in this asset class; this ID is generated off-chain

-_Metadata_: The structure containing the specific data of the asset

-_Metadata URI_: When metadata is stored off-chain, this URI indicates its storage location

## Features

### issued

Specify the asset Denom (asset category) and metadata JSON Schema to issue assets.

`CLI`

```bash
iris tx nft issue <denom> --schema=<schema-content or path/to/schema.json>
```

### Additional issuance

After the asset is issued, additional issuance (create) of specific assets of this type can be made. The asset ID, recipient address, metadata, or URI must be specified.

`CLI`

```bash
iris tx nft mint <denom> ---recipient=<recipient-address> --token-id=<token-id> --token-uri=<token-uri> --token-data=<token-data>
```

### Edit

The metadata of the specified asset can be updated.

`CLI`

```bash
iris tx nft edit <denom> <token-id> --token-uri=<token-uri> --token-data=<token-data>
```

### Transfer

Transfer designated assets.

`CLI`

```bash
iris tx nft transfer <denom> <token-id> --recipient=<recipient-address>
```

### Destroy

You can destroy the created assets.

`CLI`

```bash
iris tx nft burn <denom> <token-id>
```

### Query the specified asset class

Query asset class information based on Denom.

`CLI`

```bash
iris q query nft denom <denom>
```

### Query all asset category information

Query all issued asset class information.

`CLI`

```bash
iris q nft denoms
```

### Query the total amount of assets in a specified category

Query the total amount of assets according to Denom; accept the optional owner parameter.

`CLI`

```bash
iris q nft supply <denom> --owner=<owner>
```

### Query all assets of the specified account

Query all assets owned by an account; you can specify the Denom parameter.

`CLI`

```bash
iris q nft owner --denom=<denom>
```

### Query all assets of a specified category

Query all assets according to Denom.

`CLI`

```bash
iris q nft collection <denom>
```

### Query specified assets

Query specific assets based on Denom and ID.

`CLI`

```bash
iris q nft token <denom> <token-id>
```