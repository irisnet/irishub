# NFT

## Introduction

`NFT` provides the ability to digitize assets. Through this module, each off-chain asset will be modeled as a unique on-chain nft.

nft on the chain are identified by `ID`. With the help of the secure and non-tamperable features of the blockchain, the ownership of nft will be clarified. The transaction process of nft among members will also be publicly recorded to facilitate traceability and dispute settlement.

nft metadata (`metadata`) can be stored directly on the chain, or the `URI` of its storage source outside the chain can be stored on the chain. nft metadata is organized according to a specific [JSON Schema](https://JSON-Schema.org/). [Here](./nft-metadata.json) is an example of metadata JSON Schema.

nft need to be issued before creation to declare their abstract properties:

-_Denom_: the globally unique asset class identifier

-_Metadata Specification_: The JSON Schema that nft metadata should follow

Each specific nft is described by the following elements:

-_Denom_: the category of the nft

-_ID_: The identifier of the nft, which is unique in this nft denom; this ID is generated off-chain

-_Metadata_: The structure containing the specific data of the nft

-_Metadata URI_: When metadata is stored off-chain, this URI indicates its storage location

## Features

### issued

Specify the nft Denom (nft category) and metadata JSON Schema to issue nft.

`CLI`

```bash
iris tx nft issue <denom> --schema=<schema-content or path/to/schema.json>
```

### Additional issuance

After the nft is issued, additional issuance (create) of specific nft of this type can be made. The nft ID, recipient address, metadata, or URI must be specified.

`CLI`

```bash
iris tx nft mint <denom> ---recipient=<recipient-address> --token-id=<token-id> --token-uri=<token-uri> --token-data=<token-data>
```

### Edit

The metadata of the specified nft can be updated.

`CLI`

```bash
iris tx nft edit <denom> <token-id> --token-uri=<token-uri> --token-data=<token-data>
```

### Transfer

Transfer designated nft.

`CLI`

```bash
iris tx nft transfer <denom> <token-id> --recipient=<recipient-address>
```

### Destroy

You can destroy the created nft.

`CLI`

```bash
iris tx nft burn <denom> <token-id>
```

### Query the specified nft denom

Query nft denom information based on Denom.

`CLI`

```bash
iris q query nft denom <denom>
```

### Query all nft denom information

Query all issued nft denom information.

`CLI`

```bash
iris q nft denoms
```

### Query the total amount of nft in a specified denom

Query the total amount of nft according to Denom; accept the optional owner parameter.

`CLI`

```bash
iris q nft supply <denom> --owner=<owner>
```

### Query all nft of the specified account

Query all nft owned by an account; you can specify the Denom parameter.

`CLI`

```bash
iris q nft owner --denom=<denom>
```

### Query all nft of a specified denom

Query all nft according to Denom.

`CLI`

```bash
iris q nft collection <denom>
```

### Query specified nft

Query specific nft based on Denom and ID.

`CLI`

```bash
iris q nft token <denom> <token-id>
```