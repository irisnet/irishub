# iriscli bank account

## Description

Querying account information.

## Usage:

```
iriscli bank account [address] [flags] 
```

 

## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help     |        | False    |                       | help for account                                             |
| --chain-id     | String | False    |                       | Chain ID of tendermint node                                  |
| --height       | Int    | False    |                       | Block height to query, omit to get most recent provable block |
| --ledger       | String | False    |                       | Use a connected Ledger device                                |
| --node         | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node   | String | False    | True                  | Don't verify proofs for responses                            |



## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Query  account 

```
 iriscli bank account faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx
```

After that, you will get the detail info for the account.

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
If you query an wrrong account, you will get the fellow information.
```
 iriscli bank account faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429zz
ERROR: decoding bech32 failed: checksum failed. Expected d429yx, got d429zz.
```
If you query an empty account, you will get the fellow information.
```
iriscli bank account faa1kenrwk5k4ng70e5s9zfsttxpnlesx5ps0gfdv7
ERROR: No account with address faa1kenrwk5k4ng70e5s9zfsttxpnlesx5ps0gfdv7 was found in the state.
Are you sure there has been a transaction involving it?
```


## Extended description

Query your account in iris network.

​    



​           
