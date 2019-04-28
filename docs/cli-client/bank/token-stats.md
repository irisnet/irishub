# iriscli bank token-stats

## Description

Query the token statistic, including total loosen token, total burned token and total bonded token.

## Usage:

```
 iriscli bank token-stats <flags>
```

## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help     |        | False    |                       | Help for coin-type                                           |
| --chain-id     | String | False    |                       | Chain ID of tendermint node                                  |
| --height       | Int    | False    |                       | Block height to query, omit to get most recent provable block |
| --indent       | String | False    |                       | Add indent to JSON response                                  |
| --ledger       | String | False    |                       | Use a connected Ledger device                                |
| --node         | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node   | String | False    | True                  | Don't verify proofs for responses                            |



## Examples

### Query the token statistic

```
iriscli bank token-stats
```

Output:
```json
{
  "loosen_token": [
    "1864477.596384156921391687iris"
  ],
  "burned_token": [
    "177.59638iris"
  ],
  "bonded_token": "425182.329615843078608313iris"
}
```

​    



​           
