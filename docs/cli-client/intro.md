---
order: 1
---

# Introduction

`iriscli` is a command line client for the IRIShub network. IRIShub users can use `iriscli` to send transactions and query the blockchain data.

## Working Directory

The default working directory for the `iriscli` is `$HOME/.iriscli`, which is mainly used to save configuration files and data. The IRIShub `key` data is saved in the working directory of `iriscli`. You can also specify the `iriscli`  working directory by `--home`.

## Connecting to a Full Node

The rpc address of the `iris` node. Transactions and query requests are sent to the process listening to this port. The default is `tcp://localhost:26657`, and the rpc address can also be specified by `--node`.

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
50000
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

#### Json Indent Response

iriscli returns text format as default:

```bash
$ iriscli stake validators
Validator
  Operator Address:            iva1gfcee5u5f54kfcnufv4ypcfyldw0vu0zpwl52q
  Validator Consensus Pubkey:  icp1zcjduepquednrr0aqw4nkt8jnkhpmg4acfc7vlr0yre4uud4z0ups68hcpfsx4x9ng
  Jailed:                      false
  Status:                      Bonded
  Tokens:                      1361.0004000000246900000000000000
  Delegator Shares:            1361.0004000000246900000000000000
  Description:                 {B-2  3_C a1_}
  Unbonding Height:            0
  Minimum Unbonding Time:      1970-01-01 00:00:00 +0000 UTC
  Commission:                  rate: 0.1001000000, maxRate: 1.0000000000, maxChangeRate: 1.0000000000, updateTime: 2019-05-09 03:13:39.720700953 +0000 UTC
```

By specifing `output` and `indent`, `iriscli` can return json indent format results:

```bash
$ iriscli stake validators --output=json --indent
[
  {
    "operator_address": "iva1gfcee5u5f54kfcnufv4ypcfyldw0vu0zpwl52q",
    "consensus_pubkey": "icp1zcjduepquednrr0aqw4nkt8jnkhpmg4acfc7vlr0yre4uud4z0ups68hcpfsx4x9ng",
    "jailed": false,
    "status": 2,
    "tokens": "1361.0004000000246900000000000000",
    "delegator_shares": "1361.0004000000246900000000000000",
    "description": {
      "moniker": "B-2",
      "identity": "",
      "website": "3_C",
      "details": "a1_"
    },
    "bond_height": "0",
    "unbonding_height": "0",
    "unbonding_time": "1970-01-01T00:00:00Z",
    "commission": {
      "rate": "0.1001000000",
      "max_rate": "1.0000000000",
      "max_change_rate": "1.0000000000",
      "update_time": "2019-05-09T03:13:39.720700953Z"
    }
  }
]
```

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

Each module provides a set of command line interfaces. Here we sort these commands by modules.

| **Subcommand**                           | **Description**                                              |
| ---------------------------------------- | ------------------------------------------------------------ |
| [status](./status/README.md)             | Query remote node for status                                 |
| [tx](./tx/README.md)                     | Tx subcommands                                               |
| [tendermint](./tendermint/README.md)     | Tendermint state querying subcommands                        |
| [bank](./bank/README.md)                 | Bank subcommands for querying acccounts and sending coins etc. |
| [distribution](./distribution/README.md) | Distribution subcommands for rewards management              |
| [gov](./gov/README.md)                   | Governance and voting subcommands                            |
| [stake](./stake/README.md)               | Staking subcommands for validators and delegators            |
| [upgrade](./upgrade/README.md)           | Software Upgrade subcommands                                 |
| [service](./service/README.md)           | Service subcommands                                          |
| [guardian](./guardian/README.md)         | Guardian subcommands                                         |
| [asset](./asset.md)                      | Asset subcommands                                            |
| [rand](./rand/README.md)                 | Random Number subcommands                                    |
| [keys](./keys/README.md)                 | Keys allows you to manage your local keystore for tendermint |
| [params](./params/README.md)             | Query parameters of modules                                  |
