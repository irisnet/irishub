# iriscli bank sign

## Description

Sign transactions generated offline file. The file created with the --generate-only flag.

## Usage:

```
iriscli bank sign <file> [flags]
```

 

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | false    |                       | help for send                                                |
| --append         | boole  | True     | True                  | Append the signature to the existing ones. If disabled, old signatures would be overwritten |
| --name           | string | True     |                       | Name of private key with which to sign                       |
| --offline        | boole  | true     | false                 | Offline mode. Do not query local cache.                      |
| --print-sigs     | boole  | true     | false                 | Print the addresses that must sign the transaction and those who have already signed it, then exit |
| --chain-id       | String | False    |                       | Chain ID of tendermint node                                  |
| --account-number | int    | False    |                       | AccountNumber number to sign the tx                          |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --async          |        | false    | true                  | broadcast transactions asynchronously                        |
| --dry-run        |        | false    |                       | ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |
| --fee            | String | True     |                       | Fee to pay along with transaction                            |
| --from           | String | true     |                       | Name of private key with which to sign                       |
| --from-addr      | string | false    |                       | Specify from address in generate-only mode                   |
| --gas            | String | false    | 20000                 | gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | float  |          | 1                     | adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |
| --generate-only  |        |          |                       | build an unsigned transaction and write it to STDOUT         |
| --indent         |        |          |                       | Add indent to JSON response                                  |
| --json           |        |          |                       | return output in json format                                 |
| --memo           | String | false    |                       | Memo to send along with transaction                          |
| --print-response |        |          |                       | return tx response (only works with async = false)           |
| --sequence       | int    |          |                       | Sequence number to sign the tx                               |
| --to             | string |          |                       | Bech32 encoding address to receive coins                     |
| --ledger         | String | False    |                       | Use a connected Ledger device                                |
| --node           | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node     | String | False    | True                  | Don't verify proofs for responses                            |



## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | string   Binary encoding (hex \|b64 \|btc ) | false    | String |
| --home string         | /root/.iriscli | directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | print out full stack trace on errors        | False    |        |

## Examples

### Sign a send file 

First you must use **iriscli bank send **cli with flag **--generate-only** to generate a send recorder. Just like this.

```  
iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris --generate-only

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```



And then save the output in file  /root/output/output/node0/test_send_10iris.txt.

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
```

After that, you will get the detail info for the sign. 

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
Password to sign with 'test':
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AzlCwiA5Tvxwi7lMB/Hihfp2qnaks5Wrrgkg/Jy7sEkF"},"signature":"ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==","account_number":"0","sequence":"2"}],"memo":""}}
```





## Extended description

For offline transactions.

​    



​           
