# iriscli gov submit-proposal

## Description

Submit a proposal along with an initial deposit

## Usage

```
iriscli gov submit-proposal [flags]
```

Print help messages:

```
iriscli gov submit-proposal --help
```
## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] Deposit of proposal                                                                                                                         |          |
| --description    |                            | [string] Description of proposal                                                                                                                     | Yes      |
| --key            |                            | The key of parameter                                                                                                                                 |          |
| --op             |                            | [string] The operation of parameter                                                                                                                  |          |
| --param          |                            | [string] Parameter of proposal,eg. [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] The path of param.json                                                                                                                      |          |
| --title          |                            | [string] Title of proposal                                                                                                                           | Yes      |
| --type           |                            | [string] ProposalType of proposal,eg:Text/ParameterChange/SoftwareUpgrade                                                                            | Yes      |
| --version           |            0                | [uint64] the version of the new protocol                                                                            |       |
| --software           |           " "                 | [string] the software of the new protocol                                                                         |       |
| --switch-height           |       0                     | [string] the switchheight of the new protocol                                                         |       |

## Examples

### Submit a 'Text' type proposal

```shell
iriscli gov submit-proposal --chain-id=test --title="notice proposal" --type=Text --description="a new text proposal" --from=node0 --fee=0.01iris
```

After you enter the correct password, you're done with submitting a new proposal, and then remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt
Password to sign with 'node0':
Committed at block 44 (tx hash: 2C28A87A6262CACEDDB4EBBC60FE989D0DB2B7DEB1EC6795D2F4707DA32C7CBF, response: {Code:0 Data:[49] Log:Msg 0:  Info: GasWanted:200000 GasUsed:8091 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 112 114 111 112 111 115 97 108] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 101 114] Value:[102 97 97 49 115 108 116 106 120 100 103 107 48 48 115 56 54 50 57 50 122 48 99 110 55 97 53 100 106 99 99 116 54 101 115 115 110 97 118 100 121 122] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 97 114 97 109] Value:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 52 48 52 53 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-proposal",
     "completeConsumedTxFee-iris-atto": "\"4045500000000000\"",
     "param": "",
     "proposal-id": "1",
     "proposer": "faa1sltjxdgk00s86292z0cn7a5djcct6essnavdyz"
   }
 }
```

### Submit a 'ParameterChange' type proposal

```shell
iriscli gov submit-proposal --chain-id=test --title="update MinDeposit proposal" --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --type=ParameterChange --description="a new parameter change proposal" --from=node0 --fee=0.01iris
```

After that, you're done with submitting a new 'ParameterChange' proposal. 
The details of changed parameters （get parameters through query-params, modify it and then add "update" on the "op", more details in usage scenarios）and other fields of proposal are similar with text proposal.
Note: in this case, --path and --param cannot be both empty.

### Submit a 'SoftwareUpgrade' type proposal

```shell
iriscli gov submit-proposal --chain-id=test --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.01iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80
```

In this case, 'title'、 'type' and 'desciption' of the proposal is required parameters, also you should back up your proposal-id which is the only way to retrieve your proposal.
