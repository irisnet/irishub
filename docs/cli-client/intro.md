---
order: 1
---

# Introduction

`iriscli` is a command line client for the IRIShub network. IRIShub users can use `iriscli` to send transactions and query the blockchain data.

## Working Directory

The default working directory for the `iriscli` is `$HOME/.iriscli`, which is mainly used to save configuration files and data. The IRIShub `key` data is saved in the working directory of `iriscli`. You can also specify the `iriscli`  working directory by `--home`.

## Connecting to a Full Node

The `iris` node provides a RPC interface, transactions and query requests are sent to the process listening to it. The default rpc address the `iriscli` is connected to is `tcp://localhost:26657`, it can also be specified by `--node`.

## Setting Default Configs

The `iriscli config` command interactively configures some default parameters, such as chain-id, home, fee, and node.

E.g.

```bash
$ iriscli config
> Where is your iriscli home directory? (Default: ~/.iriscli)
/root/my_cli_home
> Where is your validator node running? (Default: tcp://localhost:26657)
tcp://192.168.0.1:26657
Do you trust this node? [y/n]:y
> What is your chainID?
irishub
> Please specify default fee
0.3iris
```

## Global Flags

### GET Commands

All GET commands has the following global flags:

| Name, shorthand | type   | Required | Default Value         | Description                                                   |
| --------------- | ----   | -------- | --------------------- | ------------------------------------------------------------- |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                   |
| --height        | int    |          | 0                     | Block height to query, omit to get most recent provable block |
| --help, -h      | string |          |                       | Print help message                                            |
| --output        | string |          | text                  | Response format text or json                                  |
| --indent        | bool   |          | false                 | Add indent to JSON response                                   |
| --ledger        | bool   |          | false                 | Use a connected Ledger device                                 |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>` to tendermint rpc interface for this chain    |
| --trust-node    | bool   |          | true                  | Don't verify proofs for responses                             |

### POST Commands

All POST commands have the following global flags:

| Name, shorthand  | type   | Required | Default               | Description                                                                                                    |
| -----------------| -----  | -------- | --------------------- | -------------------------------------------------------------------------------------------------------------- |
| --account-number | int    |          | 0                     | AccountNumber to sign the tx                                                                                   |
| --async          | bool   |          | false                 | Broadcast transactions asynchronously(only works with commit = false)                                          |
| --commit         | bool   |          | false                 | Broadcast transaction and wait until the transaction is included by a block                                    |
| --chain-id       | string | true     |                       | Chain ID of tendermint node                                                                                    |
| --dry-run        | bool   |          | false                 | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                        |
| --fee            | string | true     |                       | Fee to pay along with transaction                                                                              |
| --from           | string |          |                       | Name of private key with which to sign                                                                         |
| --from-addr      | string |          |                       | Specify from address in generate-only mode                                                                     |
| --gas            | int    |          | 50000                 | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                    |
| --gas-adjustment | int    |          | 1.5                   | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set |
| --generate-only  | bool   |          | false                 | Build an unsigned transaction and write it to STDOUT                                                           |
| --help, -h       | string |          |                       | Print help message                                                                                             |
| --indent         | bool   |          | false                 | Add indent to JSON response                                                                                    |
| --json           | string |          | false                 | Return output in json format                                                                                   |
| --ledger         | bool   |          | false                 | Use a connected Ledger device                                                                                  |
| --memo           | string |          |                       | Memo to send along with transaction                                                                            |
| --node           | string |          | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain                                                     |
| --print-response | bool   |          | false                 | Return tx response (only works with async = false)                                                             |
| --sequence       | int    |          | 0                     | Sequence number to sign the tx                                                                                 |
| --trust-node     | bool   |          | true                  | Don't verify proofs for responses                                                                              |

## Module Commands

| **Subcommand**                           | **Description**                                              |
| ---------------------------------------- | ------------------------------------------------------------ |
| [status](./status.md)             | Query remote node for status                                 |
| [tx](./tx.md)                     | Tx subcommands                                               |
| [tendermint](./tendermint.md)     | Tendermint state querying subcommands                        |
| [bank](./bank.md)                 | Bank subcommands for querying acccounts and sending coins etc. |
| [distribution](./distribution.md) | Distribution subcommands for rewards management              |
| [gov](./gov.md)                   | Governance and voting subcommands                            |
| [stake](./stake.md)               | Staking subcommands for validators and delegators            |
| [upgrade](./upgrade.md)           | Software Upgrade subcommands                                 |
| [service](./service.md)           | Service subcommands                                          |
| [guardian](./guardian.md)         | Guardian subcommands                                         |
| [asset](./asset.md)                      | Asset subcommands                                            |
| [rand](./rand.md)                 | Random Number subcommands                                    |
| [keys](./keys.md)                 | Keys allows you to manage your local keystore for tendermint |
| [params](./params.md)             | Query parameters of modules                                  |
