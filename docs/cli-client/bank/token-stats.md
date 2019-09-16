# iriscli bank token-stats

## Description

Query the token statistic, including total loose tokens, total burned token and total bonded token.

## Usage

```
 iriscli bank token-stats <tokenId> [flags]
```

## Flags

| Name, shorthand | Type   | Required | Default               | Description                                                  |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help      |        |          |                       | Help for coin-type                                           |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                  |
| --height        | int    |          |                       | Block height to query, omit to get most recent provable block|
| --indent        | string |          |                       | Add indent to JSON response                                  |
| --ledger        | string |          |                       | Use a connected Ledger device                                |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>`to tendermint rpc interface for this chain    |
| --trust-node    | string |          | true                  | Don't verify proofs for responses                            |



## Examples

### Query the token statistic

```
iriscli bank token-stats iris
```

Output:
```
TokenStats:
  Loose Token:  
    denom:iris
    amount:1864477.596384156921391687
  Burned Token:
    denom:iris
    amount:7177.596384156921391687
  Bonded Token:  
    denom:iris
    amount:1857300.596384156921391687
```

​    



​           
