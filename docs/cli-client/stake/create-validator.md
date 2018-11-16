# iriscli stake create-validator

## Description

Create new validator initialized with a self-delegation to it

## Usage

```
iriscli stake create-validator [flags]
```

## Flags

| Name, shorthand              | Default               | Description                                                         | Required |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] AccountNumber number to sign the tx                           |          |
| --address-delegator          |                       | [string] Bech address of the delegator                                       |          |
| --amount                     |                       | [string] Amount of coins to bond                                             | Yes      |
| --async                      |                       | Broadcast transactions asynchronously                               |          |
| --chain-id                   |                       | [string] Chain ID of tendermint node                                | Yes      |
| --commission-max-change-rate |                       | [string] The maximum commission change rate percentage (per day)    | Yes      |
| --commission-max-rate        |                       | [string] The maximum commission rate percentage                              | Yes      |
| --commission-rate            |                       | [string] The initial commission rate percentage                              | Yes      |
| --details                    |                       | [string] Optional details                                                    |          |
| --dry-run                    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |          |
| --fee                        |                       | [string] Fee to pay along with transaction                                   | Yes      |
| --from                       |                       | [string] Name of private key with which to sign                              | Yes      |
| --from-addr                  |                       | [string] Specify from address in generate-only mode                          |          |
| --gas                        | 200000                | [string] Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |           |
| --gas-adjustment             | 1                     | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignor |          |
| --generate-only              |                       | Build an unsigned transaction and write it to STDOUT                |          |
| --genesis-format             |                       | Export the transaction in gen-tx format; it implies --generate-only |          |
| --help, -h                   |                       | Help for create-validator                                           |          |
| --identity                   |                       | [string] Optional identity signature (ex. UPort or Keybase)         |          |
| --indent                     |                       | Add indent to JSON response                                         |          |
| --ip                         |                       | [string] Node's public IP. It takes effect only when used in combination with --genesis-format |           |
| --json                       |                       | Return output in json format                                        |          |
| --ledger                     |                       | Use a connected Ledger device                                       |          |
| --memo                       |                       | [string] Memo to send along with transaction                        |          |
| --moniker                    |                       | [string] Validator name                                             |          |
| --node                       | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --node-id                    |                       | [string] Node's ID                                                  |          |
| --print-response             |                       | Return tx response (only works with async = false)                  |          |
| --pubkey                     |                       | [string] Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 | Yes       |
| --sequence                   |                       | [int] Sequence number to sign the tx                                |          |
| --trust-node                 | true                  | Don't verify proofs for responses                                   |          |
| --website                    |                       | [string] Optional website                                                    |          |

## Examples

### Create new validator

```shell
iriscli stake create-validator --chain-id=ChainID --from=KeyName --fee=Fee --pubkey=ValidatorPublicKey --commission-max-change-rate=CommissionMaxChangeRate --commission-max-rate=CommissionMaxRate --commission-rate=CommissionRate --amount=Coins
```

After that, you're done with creating a new validator.

```txt
TODO
```
