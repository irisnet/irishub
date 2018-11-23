# iriscli gov vote

## Description

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

## Usage

```
iriscli gov vote [flags]
```

Print help messages:

```
iriscli gov vote --help
```
## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --option         |                            | [string] Vote option {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --proposal-id    |                            | [string] ProposalID of proposal voting on                                                                                                            | Yes      |

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
