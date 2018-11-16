# iriscli gov vote

## Description

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

## Usage

```
iriscli gov vote [flags]
```

## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                                                                            |          |
| --async          |                            | broadcast transactions asynchronously                                                                                                                |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
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
| --option         |                            | [string] vote option {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --print-response |                            | return tx response (only works with async = false)                                                                                                   |          |
| --proposal-id    |                            | [string] proposalID of proposal voting on                                                                                                            | Yes      |
| --sequence       |                            | [int] Sequence number to sign the tx                                                                                                                 |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Vote for proposal

```shell
iriscli gov vote --chain-id=test --proposal-id=1 --option=Yes --from node0 --fee=0.01iris
```

Validators and delegators can vote for proposals which enter voting period.
After you enter the correct password, you're done with voting for a 'VotingPeriod' proposal.

```txt
Vote[Voter:faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd,ProposalID:1,Option:Yes]Password to sign with 'node0':
Committed at block 989 (tx hash: ABDD88AA3CA8F365AC499427477CCE83ADF09D7FC2D62643D0217107E489A483, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:2408 Tags:[{Key:[97 99 116 105 111 110] Value:[118 111 116 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[118 111 116 101 114] Value:[102 97 97 49 52 113 53 114 102 57 115 108 50 100 113 100 50 117 120 114 120 121 107 97 102 120 113 51 110 117 51 108 106 50 102 112 57 108 55 112 103 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 50 48 52 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "vote",
     "completeConsumedTxFee-iris-atto": "\"120400000000000\"",
     "proposal-id": "1",
     "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd"
   }
 }
```