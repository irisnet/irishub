# iriscli bank coin-type

## Description

Query a special kind of token in IRISnet. The native token in IRIShub is `iris`, which has following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`. 

## Usage

```
 iriscli bank coin-type <coin_name> [flags]
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

### Query native token `iris`

```
iriscli bank coin-type iris
```

After that, you will get the detail info for the native token `iris`

```
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```



## Extended description

Query a special token in IRISnet.

​    



​           
