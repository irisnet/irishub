# iriscli gov deposit

## Description
 
Deposit tokens for activing proposal
 
## Usage
 
```
iriscli gov deposit [flags]
```

## Flags
 
| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                                                                            |          |
| --async          |                            | broadcast transactions asynchronously                                                                                                                |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --deposit        |                            | [string] deposit of proposal                                                                                                                         |          |
| --dry-run        |                            | ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                                                              |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                                                                           | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                                                                                      | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                                                                                  |          |
| --gas            | 200000                     | [string] gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                                                 |          |
| --gas-adjustment | 1                          | [float] adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                                                                                 |          |
| --help, -h       |                            | help for submit-proposal                                                                                                                             |          |
| --indent         |                            | Add indent to JSON response                                                                                                                          |          |
| --json           |                            | return output in json format                                                                                                                         |          |
| --ledger         |                            | Use a connected Ledger device                                                                                                                        |          |
| --memo           |                            | [string] Memo to send along with transaction                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --print-response |                            | return tx response (only works with async = false)                                                                                                   |          |
| --proposal-id    |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --sequence       |                            | [int] Sequence number to sign the tx                                                                                                                 |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Deposit

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=10iris --from=node0 --fee=0.1iris
```

 After that, you're done with depositing iris tokens for an activing proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt
Password to sign with 'node0':
Committed at block 861 (tx hash: 42D72A67ADCBE1FD90D8313E3EFB5F63A626B41F16DC0A0C7FD116907604CEF6, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:6629 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 112 111 115 105 116 101 114] Value:[102 97 97 49 115 108 116 106 120 100 103 107 48 48 115 56 54 50 57 50 122 48 99 110 55 97 53 100 106 99 99 116 54 101 115 115 110 97 118 100 121 122] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 51 51 49 52 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "deposit",
     "completeConsumedTxFee-iris-atto": "\"3314500000000000\"",
     "depositer": "faa1sltjxdgk00s86292z0cn7a5djcct6essnavdyz",
     "proposal-id": "1"
   }
 }
```