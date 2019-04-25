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
| --deposit        |                            | [string] Deposit of proposal(at least  30% of minDeposit)                                                                               |          |
| --description    |                            | [string] Description of proposal                                                                                                                     | Yes      |
| --key            |                            | The key of parameter                                                                                                                                 |          |
| --op             |                            | [string] The operation of parameter                                                                                                                  |          |
| --param          |                            | [string] Parameter of proposal,eg. [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] The path of param.json                                                                                                                      |          |
| --title          |                            | [string] Title of proposal                                                                                                                           | Yes      |
| --type           |                            | [string] ProposalType of proposal,eg:ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage                                                                            | Yes      |
| --version           |            0                | [uint64] the version of the new protocol                                                                            |       |
| --software           |           " "                 | [string] the software of the new protocol                                                                         |       |
| --switch-height           |       0                     | [string] the switchheight of the new protocol                                                         |       |
| --threshold | "0.8"   |  [string] the upgrade signal threshold of the software upgrade                                                   |               |

## Examples

The proposor which submits proposal mortgages at least 30% of `MinDeposit` ,  detail in [Gov](../../feature/governance.md)

### Submit a 'ParameterChange' type proposal

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title="update MinDeposit proposal" --param='mint/Inflation=0.050' --type=ParameterChange --description="a new parameter change proposal" --from=node0 --fee=0.3iris --commit
```

After that, you're done with submitting a new 'ParameterChange' proposal. 
The details of changed parameters （get parameters through query-params, modify it and then add "update" on the "op", more details in usage scenarios）and other fields of proposal are similar with text proposal.
Note: in this case, --path and --param cannot be both empty.

### Submit a 'SoftwareUpgrade' type proposal

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.3iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80 --threshold=0.9 --commit
```

In this case, 'title'、 'type' and 'desciption' of the proposal is required parameters, also you should back up your proposal-id which is the only way to retrieve your proposal.
