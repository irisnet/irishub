# Asset

## Introduction

This specification describes asset management on IRISHub. Anyone could issue a new asset on IRISHub, or propose pegging an existing asset from any other blockchains via On-Chain Governance.

## Concepts

### Assets

#### Native Assets

IRISHub allows individuals and companies to create and issue their own tokens for anything they can imagine. The potential use cases are innumerable. On the one hand, native assets can be used as simple event tickets deposited on the customers' mobile phones to pass the entrance of a concert. On the other hand, they can be used for crowd funding, ownership tracking or even to sell equity of a company in form of stock.
All you need to do in order to create a new `native asset` is to execute a one-line command, defining your preferred parameters for your coin, such as supply, symbol, description etc. From that point on, you can issue some of your coins to whomever you want just like sending iris.
As the owner of that asset, you don’t need to take care of any technical details of blockchain technology, such as distributed consensus algorithms, blockchain development or integration. You don’t even need to run any mining equipment or servers, at all.

#### External Assets

Instead of creating a `native asset` where the full control over supply is under the issuer, we can also create an `external asset` which already exists on another blockchain and let the market deal with demand and supply.
The only way to create an `external asset` is by submitting an `TokenAddition` proposal via Governance, except that the top 20 CMC tokens are pre-configured in the system for users' convenience.

### Gateways

A gateway is a trusted party that facilitates moving value into and out of the IRIS Network. Gateways are basically equivalent to the standard exchange model where you depend on the solvency of the exchange to be able to redeem your coins. Generally gateways issue [native assets](#native-assets) prefixed with their symbol, like GDEX, OPEN, and so on. These assets are backed 100% by the real BTC or ETH or any other coin that people deposit in the gateways.

### Fees

#### Related parameters

| name                   | Type | Default    | Description                                                                 |
| ---------------------- | ---- | ---------- | --------------------------------------------------------------------------- |
| AssetTaxRate           | Dec  | 0.4        | Asset Tax Rate, which is the ratio of Community Tax                         |
| IssueTokenBaseFee      | Coin | 60000iris  | Benchmark fees for issuing Fungible Token                                   |
| MintTokenFeeRatio      | Dec  | 0.1        | Fungible Token mint rate (relative to the issue fee)                        |
| CreateGatewayBaseFee   | Coin | 120000iris | Benchmark fees for creating Gateways                                        |
| GatewayAssetFeeRatio   | Dec  | 0.1        | Rate of issuing gateway tokens (relative to the issue fee of native tokens) |

Note: The parameters above can all be governed.

#### Fee of creating a gateway

- Baseline-Fee: The base fee required to create a gateway, ie the fee of the gateway with the minimum (3) length of `Moniker`
- Fee-Factor calculation formula: (ln(len({moniker}))/ln3)^4
- Gateway-Create-Fee calculation formula: CreateGatewayBaseFee/Fee-Factor; round up the result to iris (rounded to greater than 1 and 1 for less than or equal to 1)

#### Fee of issuing a fungible token

- Baseline-Fee: The basic fee required to issue the FT, ie the fee for the minimum FT Symbol length (3)
- Fee-Factor calculation formula: (ln(len({symbol}))/ln3)^4
- FT-Issue-Fee calculation formula: IssueFTBaseFee/Fee-Factor; round the result to iris (rounded to greater than 1 and 1 for less than or equal to 1)

#### Fee of minting a fungible token

- FT-Mint-Rate: Relative to the rate at which the FT is issued
- FT-Mint-Fee calculation formula: FT-Issue-Fee * MintFTFeeRatio; the result is rounded to iris (rounded to greater than 1 and 1 for less than or equal to 1)
  
#### Fee of issuing/minting a gateway token

- Gateway-Token-Rate (Issue/Mint): Relative to the rate at which the Native FT is issued/minted
- Gateway-Token-Fee calculation formula: (issued/minted Native FT Fee) * GatewayAssetFeeRatio; the result is rounded to iris (rounded to greater than 1 and 1 when less than or equal to 1)

#### Fee deduction method

- Community Tax: Part of the asset-related operating expenses will be used as the Community Tax, and the ratio will be determined by AssetTaxRate.
- Burned: The rest will be burned

## Actions

- **Tokens**

  - [Issue Token](../cli-client/asset.md#iriscli-asset-issue-token)

    - [Issue a native token](../cli-client/asset.md#issue-native-token)

    - [Issue a gateway token](../cli-client/asset.md#issue-gateway-token)

    - [Send tokens](../cli-client/asset.md#send-tokens)

  - [Query Token](../cli-client/asset.md#iriscli-asset-query-token)

  - [Query Tokens](../cli-client/asset.md#iriscli-asset-query-tokens)

  - [Edit Token](../cli-client/asset.md#iriscli-asset-edit-token)

  - [Mint Token](../cli-client/asset.md#iriscli-asset-mint-token)

  - [Burn Token](../cli-client/bank.md#iriscli-bank-burn)

  - [Transfer Ownership](../cli-client/asset.md#iriscli-asset-transfer-token-owner)

- **Gateways**

  - [Create Gateway](../cli-client/asset.md#iriscli-asset-create-gateway)

  - [Query Gateway](../cli-client/asset.md#iriscli-asset-query-gateway)

  - [Query Gateways](../cli-client/asset.md#iriscli-asset-query-gateways)

  - [Edit Gateway](../cli-client/asset.md#iriscli-asset-edit-gateway)

  - [Transfer Ownership](../cli-client/asset.md#iriscli-asset-transfer-gateway-owner)

- **Fees**

  - [Query Gateway Fee](../cli-client/asset.md#query-fee-of-creating-a-gateway)

  - [Query Native Token Fee](../cli-client/asset.md#query-fee-of-issuing-and-minting-a-native-token)

  - [Query Gateway Token Fee](../cli-client/asset.md#query-fee-of-issuing-and-minting-a-gateway-token)
