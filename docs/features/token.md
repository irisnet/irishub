# Asset

## Introduction

This specification describes asset management on IRISHub. Anyone could issue a new asset on IRISHub, or propose pegging an existing asset from any other blockchains via On-Chain Governance.

## Concepts

### Assets

#### Native Assets

IRISHub allows individuals and companies to create and issue their own tokens for anything they can imagine. The potential use cases are innumerable. On the one hand, native assets can be used as simple event tickets deposited on the customers' mobile phones to pass the entrance of a concert. On the other hand, they can be used for crowd funding, ownership tracking or even to sell equity of a company in form of stock.
All you need to do in order to create a new `native asset` is to execute a one-line command, defining your preferred parameters for your coin, such as supply, symbol, description etc. From that point on, you can issue some of your coins to whomever you want just like sending iris.
As the owner of that asset, you don’t need to take care of any technical details of blockchain technology, such as distributed consensus algorithms, blockchain development or integration. You don’t even need to run any mining equipment or servers, at all.

### Fees

#### Related parameters

| name                   | Type | Default    | Description                                                                 |
| ---------------------- | ---- | ---------- | --------------------------------------------------------------------------- |
| AssetTaxRate           | Dec  | 0.4        | Asset Tax Rate, which is the ratio of Community Tax                         |
| IssueTokenBaseFee      | Coin | 60000iris  | Benchmark fees for issuing Fungible Token                                   |
| MintTokenFeeRatio      | Dec  | 0.1        | Fungible Token mint rate (relative to the issue fee)                        |

Note: The parameters above can all be governed.

#### Fee of issuing a fungible token

- Baseline-Fee: The basic fee required to issue the FT, ie the fee for the minimum FT Symbol length (3)
- Fee-Factor calculation formula: (ln(len({symbol}))/ln3)^4
- FT-Issue-Fee calculation formula: IssueFTBaseFee/Fee-Factor; round the result to iris (rounded to greater than 1 and 1 for less than or equal to 1)

#### Fee of minting a fungible token

- FT-Mint-Rate: Relative to the rate at which the FT is issued
- FT-Mint-Fee calculation formula: FT-Issue-Fee * MintFTFeeRatio; the result is rounded to iris (rounded to greater than 1 and 1 for less than or equal to 1)

#### Fee deduction method

- Community Tax: Part of the asset-related operating expenses will be used as the Community Tax, and the ratio will be determined by AssetTaxRate.
- Burned: The rest will be burned

## Actions

- **Tokens**

  - [Issue Token](../cli-client/token.md#iriscli-asset-token-issue)

    - [Issue a native token](../cli-client/token.md#issue-a-token)
  
    - [Send tokens](../cli-client/token.md#send-tokens)

  - [Query Tokens](../cli-client/token.md#iriscli-asset-token-tokens)

  - [Edit Token](../cli-client/token.md#iriscli-asset-token-edit)

  - [Mint Token](../cli-client/token.md#iriscli-asset-token-mint)

  - [Burn Token](../cli-client/bank.md#iriscli-bank-burn)

  - [Transfer Ownership](../cli-client/token.md#iriscli-asset-token-transfer)

- **Fees**
  
  - [Query Native Token Fee](../cli-client/token.md#query-fee-of-issuing-and-minting-a-token)
  