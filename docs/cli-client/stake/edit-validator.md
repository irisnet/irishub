# iriscli stake edit-validator

## Description

Edit existing validator account

## Usage

```
iriscli stake edit-validator [flags]
```

## Flags

| Name, shorthand              | Default               | Description                                                         | Required |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] AccountNumber number to sign the tx                           |          |
| --address-delegator          |                       | [string] Bech address of the delegator                                       |          |
| --amount                     |                       | [string] Amount of coins to bond                                             |          |
| --async                      |                       | Broadcast transactions asynchronously                               |          |
| --chain-id                   |                       | [string] Chain ID of tendermint node                                | Yes      |
| --commission-max-change-rate |                       | [string] The maximum commission change rate percentage (per day)    |          |
| --commission-max-rate        |                       | [string] The maximum commission rate percentage                              |          |
| --commission-rate            |                       | [string] The initial commission rate percentage                              |          |
| --details                    |                       | [string] Optional details                                                    |          |
| --dry-run                    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it |          |
| --fee                        |                       | [string] Fee to pay along with transaction                                   | Yes      |
| --from                       |                       | [string] Name of private key with which to sign                              | Yes      |
| --from-addr                  |                       | [string] Specify from address in generate-only mode                          |          |
| --gas                        | 200000                | [string] Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |           |
| --gas-adjustment             | 1                     | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignor |          |
| --generate-only              |                       | Build an unsigned transaction and write it to STDOUT                |          |
| --genesis-format             |                       | Export the transaction in gen-tx format; it implies --generate-only |          |
| --help, -h                   |                       | Help for edit-validator                                           |          |
| --identity                   |                       | [string] Optional identity signature (ex. UPort or Keybase)         |          |
| --indent                     |                       | Add indent to JSON response                                         |          |
| --ip                         |                       | [string] Node's public IP. It takes effect only when used in combination with --genesis-format |           |
| --json                       |                       | Return output in json format                                        |          |
| --ledger                     |                       | Use a connected Ledger device                                       |          |
| --memo                       |                       | [string] Memo to send along with transaction                        |          |
| --moniker                    |                       | [string] Validator name                                             |          |
| --node                       | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --node-id                    |                       | [string] Node's ID                                                  |          |
| --print-response             |                       | Return tx response (only works with async = false)                  |          |
| --pubkey                     |                       | [string] Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |           |
| --sequence                   |                       | [int] Sequence number to sign the tx                                |          |
| --trust-node                 | true                  | Don't verify proofs for responses                                   |          |
| --website                    |                       | [string] Optional website                                                    |          |

## Examples

### Edit existing validator account

```shell
iriscli stake edit-validator --from=KeyName --chain-id=ChainID --fee=Fee --memo=YourMemo
```

After that, you're done with editting a new validator.

```txt
Committed at block 2160 (tx hash: C48CABDA1183B5319003433EB1FDEBE5A626E00BD319F1A84D84B6247E9224D1, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3540 Tags:[{Key:[97 99 116 105 111 110] Value:[101 100 105 116 45 118 97 108 105 100 97 116 111 114] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[109 111 110 105 107 101 114] Value:[117 98 117 110 116 117 49 56] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[105 100 101 110 116 105 116 121] Value:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 55 55 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "edit-validator",
     "completeConsumedTxFee-iris-atto": "\"177000000000000\"",
     "destination-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
     "identity": "",
     "moniker": "ubuntu18"
   }
}
```
