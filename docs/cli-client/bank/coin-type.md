# iriscli bank coin-type

## Description

Query  a special  kind of token in IRISnet. The native token in IRIShub is `iris`, which has following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`. 

## Usage:

```
 iriscli bank coin-type [coin_name] [flags]
```

 

## Examples

```     
iriscli bank coin-type iris
```



## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help     |        | false    |                       | help for coin-type                                           |
| --chain-id     | String | False    |                       | Chain ID of tendermint node                                  |
| --height       | Int    | False    |                       | block height to query, omit to get most recent provable block |
| --indent       | String | False    |                       | Add indent to JSON response                                  |
| --ledger       | String | False    |                       | Use a connected Ledger device                                |
| --node         | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node   | String | False    | True                  | Don't verify proofs for responses                            |



## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | string   Binary encoding (hex \|b64 \|btc ) | false    | String |
| --home string         | /root/.iriscli | directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | print out full stack trace on errors        | False    |        |

## Extended description

Query a special token in iris network.

​    



​           
