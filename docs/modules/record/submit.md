# iriscli record submit

## Description

Submit a record on chain

## Usage

```
iriscli record submit [flags]
```

## Flags

| Name, shorthand  | Default                    | Description                                                                                 | Required |
| ---------------  | -------------------------- | ------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                   |          |
| --async          |                            | Broadcast transactions asynchronously                                                       |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                        | Yes      |
| --description    | description                | [string] Uploaded file description                                                          |          |
| --dry-run        |                            | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it     |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                  | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                             | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                         |          |
| --gas string     | 200000                     | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |          |
| --gas-adjustment | 1                          | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                        |          |
| -h, --help       |                            | help for submit                                                                             |          |
| --indent         |                            | Add indent to JSON response                                                                 |          |
| --json           |                            | return output in json format                                                                |          |
| --ledger         |                            | Use a connected Ledger device                                                               |          |
| --memo           |                            | [string] Memo to send along with transaction                                                |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                           |          |
| --onchain-data   |                            | [string] on chain data source                                                               | Yes      |
| --print-response |                            | return tx response (only works with async = false)                                          |          |
| --sequence       |                            | [int] Sequence number to sign the tx                                                        |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                           |          |

## Examples

### Submit a record

```shell
iriscli record submit --chain-id="test" --onchain-data="this is my on chain data" --from=node0 --fee=0.1iris
```

After that, you're done with submitting a new record, but remember to back up your record id, it's the only way to retrieve your on chain record.

```txt
Password to sign with 'node0':
Committed at block 72 (tx hash: 7CCC8B4018D4447E6A496923944870E350A1A3AF9E15DB15B8943DAD7B5D782B, response: {Code:0 Data:[114 101 99 111 114 100 58 97 98 53 54 48 50 98 97 99 49 51 102 49 49 55 51 55 101 56 55 57 56 100 100 53 55 56 54 57 99 52 54 56 49 57 52 101 102 97 100 50 100 98 51 55 54 50 53 55 57 53 102 49 101 102 100 56 100 57 100 54 51 99 54] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4090 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 114 101 99 111 114 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[111 119 110 101 114 65 100 100 114 101 115 115] Value:[102 97 97 49 50 50 117 122 122 112 117 103 116 114 122 115 48 57 110 102 51 117 104 56 120 102 106 97 122 97 53 57 120 118 102 57 114 118 116 104 100 108] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 111 114 100 45 105 100] Value:[114 101 99 111 114 100 58 97 98 53 54 48 50 98 97 99 49 51 102 49 49 55 51 55 101 56 55 57 56 100 100 53 55 56 54 57 99 52 54 56 49 57 52 101 102 97 100 50 100 98 51 55 54 50 53 55 57 53 102 49 101 102 100 56 100 57 100 54 51 99 54] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 50 48 52 53 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-record",
     "completeConsumedTxFee-iris-atto": "\"2045000000000000\"",
     "ownerAddress": "faa122uzzpugtrzs09nf3uh8xfjaza59xvf9rvthdl",
     "record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
   }
 }
```

