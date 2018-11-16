# iriscli stake unjail

## Description

Unjail validator previously jailed for downtime

## Usage

```
iriscli stake redelegate [flags]
```

## Flags

| Name, shorthand              | Default               | Description                                                         | Required |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] AccountNumber number to sign the tx                           |          |
| --address-validator-dest     |                       | [string] Bech address of the destination validator                  |          |
| --address-validator-source   |                       | [string] Bech address of the source validator                       |          |
| --async                      |                       | Broadcast transactions asynchronously                               |          |
| --chain-id                   |                       | [string] Chain ID of tendermint node                                | Yes      |
| --dry-run                    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it ||
| --fee                        |                       | [string] Fee to pay along with transaction                          | Yes      |
| --from                       |                       | [string] Name of private key with which to sign                     | Yes      |
| --from-addr                  |                       | [string] Specify from address in generate-only mode                 |          |
| --gas                        | 200000                | [string] Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically ||
| --gas-adjustment             | 1                     | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ||
| --generate-only              |                       | build an unsigned transaction and write it to STDOUT                |          |
| --help, -h                   |                       | help for unjail                                                 |          |
| --indent                     |                       | Add indent to JSON response                                         |          |
| --json                       |                       | return output in json format                                        |          |
| --ledger                     |                       | Use a connected Ledger device                                       |          |
| --memo                       |                       | [string] Memo to send along with transaction                                 |          |
| --node                       | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain          |          |
| --print-response             |                       | return tx response (only works with async = false)                  |          |
| --sequence int               |                       | Sequence number to sign the tx                                      |          |
| --shares-amount              |                       | [string] Amount of source-shares to either unbond or redelegate as a positive integer or decimal ||
| --shares-percent             |                       | [string] Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 ||
| --trust-node                 | true                  | Don't verify proofs for responses                                   |          |

## Examples

### Unjail validator previously jailed for downtime

```shell
iriscli stake unjail --from=KeyName --fee=Fee --chain-id=ChainID
```

After that, you're done with unjailing specified validator.

