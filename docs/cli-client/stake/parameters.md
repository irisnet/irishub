# iriscli stake parameters

## Description

Query the current staking parameters information

## Usage

```
iriscli stake parameters [flags]
```

## Flags

| Name, shorthand            | Default                    | Description                                                         | Required |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id                 |                            | [string] Chain ID of tendermint node                                |          |
| --height                   | most recent provable block | block height to query                                               |          |
| --help, -h                 |                            | help for validator                                                  |          |
| --indent                   |                            | Add indent to JSON response                                         |          |
| --ledger                   |                            | Use a connected Ledger device                                       |          |
| --node                     | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node               | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query the current staking parameters information

```shell
iriscli stake parameters
```

After that, you will get the current staking parameters information.

```txt
Params
Unbonding Time: 10m0s
Max Validators: 100:
Bonded Coin Denomination: iris-atto
```
