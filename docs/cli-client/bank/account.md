# iriscli bank account

## Description

Query  a special  account detial. 

## Usage:

```
iriscli bank account [address] [flags] 
```

 

## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help     |        | false    |                       | help for coin-type                                           |
| --chain-id     | String | False    |                       | Chain ID of tendermint node                                  |
| --height       | Int    | False    |                       | block height to query, omit to get most recent provable block |
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

## Examples

### Query  account 

```
 iriscli bank account faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx
```

After that, you will get the detail info for the special account

```
{

  "address": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx",

  "coins": [

    "50iris"

  ],

  "public_key": {

    "type": "tendermint/PubKeySecp256k1",

    "value": "AzlCwiA5Tvxwi7lMB/Hihfp2qnaks5Wrrgkg/Jy7sEkF"

  },

  "account_number": "0",

  "sequence": "1"

}



```





## Extended description

Query your account in iris network.

​    



​           
