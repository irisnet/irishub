# iriscli tx sign

## Description

Sign transactions in generated offline file. The file created with the --generate-only flag.

## Usage:

```
iriscli tx sign <file> <flags>
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
iriscli tx sign tx.json --name=<key_name> 

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}
```


After that, you will get the detail info for the sign. Like the follow output you will see the signature 

**sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==**

After signing a transaction, it could be broadcast to the network with [broadcast command](broadcast.md)