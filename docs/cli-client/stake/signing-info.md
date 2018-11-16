# iriscli stake signing-info

## Description

Query a validator's signing information

## Usage

```
iriscli stake signing-info [validator-pubkey] [flags]
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

### Query specified validator's signing information

```shell
iriscli stake signing-info ValidatorPublicKey
```

After that, you will get specified validator's signing information.

```txt
Start height: 0, index offset: 2136, jailed until: 1970-01-01 00:00:00 +0000 UTC, missed blocks counter: 0
```
