# iriscli bank send

## Description

Send token to another address. 

## Usage:

```
iriscli bank send --to=<account address> --from <key name> --fee=0.004iris --chain-id=<chain-id> --amount=10iris
```

 

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | false    |                       | help for send                                                |
| --chain-id       | String | False    |                       | Chain ID of tendermint node                                  |
| --account-number | int    | False    |                       | AccountNumber number to sign the tx                          |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --async          |        |          | true                  | broadcast transactions asynchronously                        |
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

### Send token to a address 

```
 iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris
```

After that, you will get the detail info for the send

```
[root@ce7da33d46c3 iriscli]# iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris
Password to sign with 'test':
Committed at block 2265 (tx hash: A60224C8433487D48C8B03B51CB7A2BCB014932A97A55D946E5F30E561E1195E, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4690 Tags:[{Key:[115 101 110 100 101 114] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 105 112 105 101 110 116] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 51 56 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "completeConsumedTxFee-iris-atto": "\"93800000000000\"",
     "recipient": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx",
     "sender": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx"
   }
 }

```





## Extended description

You send token to an other address.

​    



​           
