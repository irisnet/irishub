# MT

## Introduction

`MT` provides ERC-1155 capacity. The MT module is able to model off-chain assets into unique on-chain assets.

On-chain assets are identified by`ID` . With the secure and tamper-proof nature of blockchains, the ownership of assets can be confirmed and verified. Transactions of assets between multiple parties can also be publicly recorded, in order to facilitate traceability and dispute resolution. The metadata (`data` ) of assets can be stored directly on the chain, or be used to record the off-chain storage addresses.

Assets need to be issued before they are created to declare their abstract properties.

* *DenomID*: the globally unique asset class identifier, which is generated on the chain.
* DenomName: the name of the asset class.

Each specific asset is described by the following elements:

* *DenomID: the class of the asset.*
* *ID*: the identifier of the asset, which is unique within the corresponding asset class. This ID is also generated on the chain.
* Metadata: a structure that contains specific data of the assets.

## Function

### Issuance

An asset class can be created by specifying the DenomName and the creator.

`CLI`

```plain
iris tx mt issue --name=<denom-name> --from=<sender-address> --chain-id=<chain-id> --fees=<fee>
```

### Production

After issuing an asset class, a specific asset of that class can be created, during which the DenomID, number of issuance, metadata, and the addresses of both the issuer (owner of the Denom) and the receiver need to be specified.

`CLI`

```plain
iris tx mt mint <denom-id> --amount=<amount> --data=<data> --from=<sender-address> --recipient=<recipient-address> --chain-id=<chain-id> --fees=<fee>
```

### Increase of Issuance

After issuing a specific asset, the owner of the asset class can also choose to increase the issuance, which requires specifying the DenomID, the number of additional issuances, and the addresses of both the issuer (owner of the Denom) and the receiver.

`CLI`

```plain
iris tx mt mint <denom-id> --mt-id=<mt-id> --amount=<amount> --from=<sender-address> --recipient=<recipient-address> --chain-id=<chain-id> --fees=<fee>
```

### Editing

Updates can be made to the metadata of a specified asset.

`CLI`

```plain
iris tx mt edit <denom-id> <mt-id> --data=<data> --from=<sender-address> --chain-id=<chain-id> --fees=<fee>
```

### Transfer


Assets can be transferred. The amount of assets to be transferred can be specified.

`CLI`

```plain
iris tx mt transfer <sender> <recipient> <denom-id> <mt-id> <amount> --chain-id=<chain-id> --fees=<fee>
```

### Burn

Assets can be burned. The amount of assets to be burned can be specified.

`CLI`

```plain
iris tx mt burn <denom-id> <mt-id> <amount> --from=<sender-address> --chain-id=<chain-id> --fees=<fee>
```

### Query a specified asset class

Query the asset class through the DenomID.

`CLI`

```plain
iris query mt denom <denom-id>
```

### Query all asset classes

Query all issued asset classes.

`CLI`

```plain
iris query mt denoms
```

### Query the total amount of assets in a specified asset class

Query the total amount of assets through the DenomID.

`CLI`

```plain
iris query mt supply <denom-id> <mt-id>
```

### Query all assets in a specified account

Query all assets owned by an account in a specified asset class.

`CLI`

```plain
iris query mt balances <owner> <denom-id>
```

### Query specified assets

Query the information of a specific asset through the DenomID and MtID.

`CLI`

```plain
iris query mt token <denom-id> <mt-id>
```

### Query all assets in a specified asset class

Query all assets in a specified asset class through the DenomID.

`CLI`

```plain
iris query mt tokens <denom-id>
```
