# iriscli bank coin-type

## Description

Query a special kind of token in IRISnet. The native token in IRIShub is `iris`, which has following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`. 


## Usage:

```
 iriscli bank coin-type <coin_name> <flags>
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

### Query native token `iris`

```
iriscli bank coin-type iris
```

After that, you will get the detail info for the native token `iris`
```
{
  "name": "iris",
  "min_unit": {
    "denom": "iris-atto",
    "decimal": "18"
  },
  "units": [
    {
      "denom": "iris",
      "decimal": "0"
    },
    {
      "denom": "iris-milli",
      "decimal": "3"
    },
    {
      "denom": "iris-micro",
      "decimal": "6"
    },
    {
      "denom": "iris-nano",
      "decimal": "9"
    },
    {
      "denom": "iris-pico",
      "decimal": "12"
    },
    {
      "denom": "iris-femto",
      "decimal": "15"
    },
    {
      "denom": "iris-atto",
      "decimal": "18"
    }
  ],
  "origin": 1,
  "desc": "IRIS Network"
}
```



## Extended description

Query a special token in IRISnet.

​    



​           
