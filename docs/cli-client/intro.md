---
order: 1
---

# Introduction

`iris` is a command line client for the IRIShub network. IRIShub users can use `iris` to send transactions and query the blockchain data.

## Working Directory

The default working directory for the `iris` is `$HOME/.iris`, which is mainly used to save configuration files and data. The IRIShub `key` data is saved in the working directory of `iris`. You can also specify the `iris`  working directory by `--home`.

## Connecting to a Full Node

The `iris` node provides a RPC interface, transactions and query requests are sent to the process listening to it. The default rpc address the `iris` is connected to is `tcp://localhost:26657`, it can also be specified by `--node`.

## Global Flags

### GET Commands

All GET commands has the following global flags:

| Name, shorthand | type   | Required | Default Value        | Description                          |
| --------------- | ------ | -------- | -------------------- | ------------------------------------ |
| --chain-id      | string |          |                      | Chain ID of tendermint node          |
| --home          | string |          | ~/.iris | Directory for config and data        |
| --trace         | string |          |                      | Print out full stack trace on errors |

### POST Commands

All POST commands have the following global flags:

| Name, shorthand   | type   | Required | Default               | Description                                                                                                    |
| ----------------- | ------ | -------- | --------------------- | -------------------------------------------------------------------------------------------------------------- |
| --account-number  | int    |          | 0                     | AccountNumber to sign the tx                                                                                   |
| --broadcast-mode  | string |          | sync                  | Transaction broadcasting mode (sync \| async \| block)                                                         |
| --dry-run         | bool   |          | false                 | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                        |
| --fees            | string |          |                       | Fees to pay along with transaction                                                                             |
| --from            | string |          |                       | Name of private key with which to sign                                                                         |
| --gas             | string |          | 50000                 | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                    |
| --gas-adjustment  | float  |          | 1.5                   | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set |
| --gas-prices      | string |          |                       | Gas prices in decimal format to determine the transaction fee                                                  |
| --generate-only   | bool   |          | false                 | Build an unsigned transaction and write it to STDOUT                                                           |
| --help, -h        | string |          |                       | Print help message                                                                                             |
| --keyring-backend | string |          | os                    | Select keyring's backend                                                                                       |
| --ledger          | bool   |          | false                 | Use a connected Ledger device                                                                                  |
| --memo            | string |          |                       | Memo to send along with transaction                                                                            |
| --node            | string |          | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain                                                     |
| --offline         | string |          |                       | Offline mode (does not allow any online functionality)                                                         |
| --sequence        | int    |          | 0                     | Sequence number to sign the tx                                                                                 |
| --sign-mode       | string |          |                       | Choose sign mode (direct \| amino-json), this is an advanced feature                                           |
| --trust-node      | bool   |          | true                  | Don't verify proofs for responses                                                                              |
| --yes             | bool   |          | true                  | Skip tx broadcasting prompt confirmation                                                                       |
| --chain-id        | string |          |                       | Chain ID of tendermint node                                                                                    |
| --home            | string |          |                       | Directory for config and data (default "~/.iris")                                                 |
| --trace           | string |          |                       | Print out full stack trace on errors                                                                           |

## Module Commands

| **Subcommand**                    | **Description**                                                |
| --------------------------------- | -------------------------------------------------------------- |
| [bank](./bank.md)                 | Bank subcommands for querying acccounts and sending coins etc. |
| [debug](./debug.md)               | Debug subcommands                                              |
| [distribution](./distribution.md) | Distribution subcommands for rewards management                |
| [gov](./gov.md)                   | Governance and voting subcommands                              |
| [htlc](./htlc.md)                 | HTLC transaction subcommands                                   |
| [keys](./keys.md)                 | Keys allows you to manage your local keystore for tendermint   |
| [nft](./nft.md)                   | NFT subcommands                                                |
| [oracle](./oracle.md)             | Oracle transaction subcommands                                 |
| [params](./params.md)             | Query parameters of modules                                    |
| [random](./rand.md)               | Random number subcommands                                      |
| [record](./record.md)             | Record subcommands                                             |
| [slashing](./slashing.md)         | Slashing subcommands                                           |
| [service](./service.md)           | Service subcommands                                            |
| [staking](./staking.md)           | Staking subcommands for validators and delegators              |
| [status](./status.md)             | Query remote node for status                                   |
| [tendermint](./tendermint.md)     | Tendermint state querying subcommands                          |
| [token](./token.md)               | Token subcommands                                              |
| [tx](./tx.md)                     | Tx subcommands                                                 |
| [upgrade](./upgrade.md)           | Software Upgrade subcommands                                   |
