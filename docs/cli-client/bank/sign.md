# iriscli bank sign

## Description

Sign transactions in generated offline file. The file created with the --generate-only flag.

## Usage:

```
iriscli bank sign <file> [flags]
```

 

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | False    |                       | Help for send                                                |
| --append         | Boole  | True     | True                  | Append the signature to the existing ones. If disabled, old signatures would be overwritten |
| --name           | String | True     |                       | Name of private key with which to sign                       |
| --offline        | Boole  | True     | False                 | Offline mode. Do not query local cache.                      |
| --print-sigs     | Boole  | True     | False                 | Print the addresses that must sign the transaction and those who have already signed it, then exit |
| --chain-id       | String | False    |                       | Chain ID of tendermint node                                  |
| --account-number | Int    | False    |                       | AccountNumber number to sign the tx                          |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --async          |        | False    | True                  | Broadcast transactions asynchronously                        |
| --dry-run        |        | False    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |
| --fee            | String | True     |                       | Fee to pay along with transaction                            |
| --from           | String | True     |                       | Name of private key with which to sign                       |
| --from-addr      | String | False    |                       | Specify from address in generate-only mode                   |
| --gas            | String | False    | 20000                 | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | Float  | False    | 1                     | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |
| --generate-only  |        | False    |                       | Build an unsigned transaction and write it to STDOUT         |
| --indent         |        | False    |                       | Add indent to JSON response                                  |
| --json           |        | False    |                       | Return output in json format                                 |
| --memo           | String | False    |                       | Memo to send along with transaction                          |
| --print-response |        | False    |                       | Return tx response (only works with async = false)           |
| --sequence       | Int    | False    |                       | Sequence number to sign the tx                               |
| --to             | String | False    |                       | Bech32 encoding address to receive coins                     |
| --ledger         | String | False    |                       | Use a connected Ledger device                                |
| --node           | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node     | String | False    | True                  | Don't verify proofs for responses                            |



## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Sign a send file 

First you must use **iriscli bank send**  command with flag **--generate-only** to generate a send recorder. Just like this.

```  
iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris --generate-only

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```



And then save the output in file  /root/output/output/node0/test_send_10iris.txt.

Then you can sign the offline file.

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
```

After that, you will get the detail info for the sign. Like the follow output you will see the signature 

**ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==**

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
Password to sign with 'test':
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AzlCwiA5Tvxwi7lMB/Hihfp2qnaks5Wrrgkg/Jy7sEkF"},"signature":"ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==","account_number":"0","sequence":"2"}],"memo":""}}
```

After signing a transaction, it could be broadcast to the network with [broadcastc command](./broadcast.md)