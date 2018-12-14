# iriscli bank send

## Description

Sending tokens to another address. 

## Usage:

```
iriscli bank send --to=<account address> --from <key name> --fee=0.004iris --chain-id=<chain-id> --amount=10iris
```

 

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | False    |                       | Help for send                                                |
| --chain-id       | String | False    |                       | Chain ID of tendermint node                                  |
| --account-number | int    | False    |                       | AccountNumber number to sign the tx                          |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --async          |        |          | True                  | Broadcast transactions asynchronously                        |
| --dry-run        |        | False    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |
| --fee            | String | True     |                       | Fee to pay along with transaction                            |
| --from           | String | True     |                       | Name of private key with which to sign                       |
| --from-addr      | string | False    |                       | Specify from address in generate-only mode                   |
| --gas            | String | False    | 20000                 | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | Float  |          | 1                     | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |
| --generate-only  |        |          |                       | Build an unsigned transaction and write it to STDOUT         |
| --indent         |        |          |                       | Add indent to JSON response                                  |
| --json           |        |          |                       | Return output in json format                                 |
| --memo           | String | False    |                       | Memo to send along with transaction                          |
| --print-response |        |          |                       | Return tx response (only works with async = false)           |
| --sequence       | Int    |          |                       | Sequence number to sign the tx                               |
| --to             | String |          |                       | Bech32 encoding address to receive coins                     |
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

### Send token to a address 

```
 iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris
```

After that, you will get the detail info for the send

```
[Committed at block 87 (tx hash: AEA8E49C1BC9A81CAFEE8ACA3D0D96DA7B5DC43B44C06BACEC7DCA2F9C4D89FC, response:
  {
    "code": 0,
    "data": null,
    "log": "Msg 0: ",
    "info": "",
    "gas_wanted": 200000,
    "gas_used": 3839,
    "codespace": "",
    "tags": {
      "action": "send",
      "recipient": "faa1893x4l2rdshytfzvfpduecpswz7qtpstpr9x4h",
      "sender": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2"
    }
  })
```
### Common Issues

* Wrong password

```$xslt
ERROR: Ciphertext decryption failed
```
