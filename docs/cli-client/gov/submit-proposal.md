# iriscli gov submit-proposal

## Description

Submit a proposal along with an initial deposit

## Usage

```
iriscli gov submit-proposal <flags>
```

Print help messages:

```
iriscli gov submit-proposal --help
```

## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | Deposit of proposal(at least  30% of minDeposit)                                                                               |          |
| --description    |                            | Description of proposal                                                                                                                     | Yes      |
| --param          |                            | Parameter of proposal,eg. [{key:key,value:value,op:update}]                                                                                 |          |
| --title          |                            | Title of proposal                                                                                                                           | Yes      |
| --type           |                            | ProposalType of proposal,eg:ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage                                                                            | Yes      |
| --version           |            0                | the version of the new protocol                                                                            |       |
| --software           |           " "                 | the software of the new protocol                                                                         |       |
| --switch-height           |       0                     | the switch height of the new protocol                                                         |       |
| --threshold | "0.8"   |  the upgrade signal threshold of the software upgrade                                                   |               |

## Examples

The proposer should deposit at least 30% of `MinDeposit` to submit a proposal,  detailed in [Gov](../../features/governance.md)

### Submit a 'ParameterChange' type proposal

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title=<proposal_title> --param='mint/Inflation=0.050' --type=ParameterChange --description=<proposal_description> --from=<key_name> --fee=0.3iris --deposit="3000iris" 
```

Note: in this case, --path and --param cannot be both empty.

### Submit a 'SoftwareUpgrade' type proposal

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title=<proposal_title> --description=<proposal_description>  --type=SoftwareUpgrade --description=<proposal_description> --from=<key_name> --fee=0.3iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80 --threshold=0.9 --deposit="3000iris" 
```

In this case, 'title'、 'type' and 'description' of the proposal is required parameters, also you should back up your proposal-id which is the only way to retrieve your proposal.

###  How to query proposal

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
