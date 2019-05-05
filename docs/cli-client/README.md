# Command Line Client

## Global flags of query commands

All query commands has these global flags. Their unique flags will be introduced later.

| Name, shorthand | type   | Required | Default Value         | Description                                                          |
| --------------- | ----   | -------- | --------------------- | -------------------------------------------------------------------- |
| --chain-id      | string | false    | ""                    | Chain ID of tendermint node |
| --height        | int    | false    | 0                     | Block height to query, omit to get most recent provable block |
| --help, -h      | string | false    |                       | Print help message |
| --indent        | bool   | false    | false                 | Add indent to JSON response |
| --ledger        | bool   | false    | false                 | Use a connected Ledger device |
| --node          | string | false    | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain |
| --trust-node    | bool   | false    | true                  | Don't verify proofs for responses |

## Global flags of commands to send transactions

All commands which can be used to send transactions have these global flags. Their unique flags will be introduced later.

| Name, shorthand  | type   | Required | Default               | Description                                                         |
| -----------------| -----  | -------- | --------------------- | ------------------------------------------------------------------- |
| --account-number | int    | false    | 0                     | AccountNumber to sign the tx |
| --async          | bool   | false    | false                 | broadcast transactions asynchronously(only works with commit = false) |
| --commit         | bool   | false    | false                 | broadcast transaction and wait until the transaction is included by a block |
| --chain-id       | string | true     | ""                    | Chain ID of tendermint node  |
| --dry-run        | bool   | false    | false                 | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |
| --fee            | string | true     | ""                    | Fee to pay along with transaction |
| --from           | string | false    | ""                    | Name of private key with which to sign |
| --from-addr      | string | false    | ""                    | Specify from address in generate-only mode |
| --gas            | int    | false    | 50000                | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | int    | false    | 1.5                   | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set |
| --generate-only  | bool   | false    | false                 | Build an unsigned transaction and write it to STDOUT |
| --help, -h       | string | false    |                       | Print help message |
| --indent         | bool   | false    | false                 | Add indent to JSON response |
| --json           | string | false    | false                 | Return output in json format |
| --ledger         | bool   | false    | false                 | Use a connected Ledger device |
| --memo           | string | false    | ""                    | Memo to send along with transaction |
| --node           | string | false    | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain |
| --print-response | bool   | false    | false                 | return tx response (only works with async = false)|
| --sequence       | int    | false    | 0                     | Sequence number to sign the tx |
| --trust-node     | bool   | false    | true                  | Don't verify proofs for responses | 

## Module command list

Each module provides a set of command line interfaces. Here we sort these commands by modules.

1. [status command](./status/README.md)
2. [tendermint command](./tendermint/README.md)
3. [keys command](./keys/README.md)
4. [bank command](./bank/README.md)
5. [stake command](./stake/README.md)
6. [distribution command](./distribution/README.md)
7. [gov command](./gov/README.md)
8. [upgrade command](./upgrade/README.md)
9. [service command](./service/README.md)

## Config command

The `iriscli config` command interactively configures some default parameters, such as chain-id, home, fee, and node.

Example：

```
root@ubuntu16:~# iriscli config
> Where is your iriscli home directory? (Default: ~/.iriscli)
/root/my_cli_home
> Where is your validator node running? (Default: tcp://localhost:26657)
tcp://192.168.0.1:26657
Do you trust this node? [y/n]:y
> What is your chainID?
irishub
> Please specify default fee
50000

root@ubuntu16:~# iriscli status --home=/root/my_cli_home
```
