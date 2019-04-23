# iriscli bank sign

## Description

Sign transactions in generated offline file. The file created with the --generate-only flag.

## Usage:

```
iriscli bank sign <file> <flags>
```

## 标志

## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | false     |                       |  Help for sign                                        |
| --append         | Bool  | true       | True                  | Attach a signature to an existing signature.               |
| --name           | String | true       |                       | Key name for signature                                      |
| --offline        | Boole | true       | False                 | Offline mode.                                    |
| --print-sigs     | Bool  | true       | False                 | Print the address where the transaction must be signed and the signed address, then exit  |
| --chain-id       | string | true     | ""                    | Chain ID of tendermint node  |
| --dry-run        | bool   | false    | false                 | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |
| --fee            | string | true     | ""                    | Fee to pay along with transaction |
| --from           | string | false    | ""                    | Name of private key with which to sign |
| --from-addr      | string | false    | ""                    | Specify from address in generate-only mode |
| --gas            | int    | false    | 50000                | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | int    | false    | 1.5                   | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set |
| --generate-only  | bool   | false    | false                 | Build an unsigned transaction and write it to STDOUT |
| --help, -h       | string | false    |                       | Print help message |
| --indent         | bool   | false    | false                 | Add indent to JSON response |
| --json           | string | false    | false                 | Return output in json format |
| --ledger         | bool   | false    | false                 | Use a connected Ledger device |
| --memo           | string | false    | ""                    | Memo to send along with transaction |
| --node           | string | false    | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain |
| --print-response | bool   | false    | false                 | return tx response (only works with async = false)|
| --sequence       | int    | false    | 0                     | Sequence number to sign the tx |
| --trust-node     | bool   | false    | true                  | Don't verify proofs for responses | 

## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding  | hex            | String   Binary encoding (hex\b64\btc ) | False    | String |
| --home          | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text\json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Sign a send file 

First you must use `iriscli bank send` command with flag `generate-only` to generate a send recorder. Just like this.

```  
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris --generate-only

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```


And then save the output in file `tx.json`.
```
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris --generate-only > tx.json
```

Then you can sign the offline file.
```
iriscli bank sign tx.json --name=<key_name> 

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}
```


After that, you will get the detail info for the sign. Like the follow output you will see the signature 

**sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==**

After signing a transaction, it could be broadcast to the network with [broadcast command](./broadcast.md)