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
| --type           |                            | [string] ProposalType of proposal,eg:Text/ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage                                                                            | Yes      |
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
Committed at block 13 (tx hash: 234463E89B5641F9271113D72B28CA088F641DD8A63DB57257B7CAF90ED5A1C3, response:
 {
   "code": 0,
   "data": "MQ==",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 6608,
   "codespace": "",
   "tags": {
     "action": "submit_proposal",
     "param": "",
     "proposal-id": "1",
     "proposer": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz"
   }
 })
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
