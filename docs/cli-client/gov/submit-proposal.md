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
| --param          |                            | Parameter of proposal,eg. mint/Inflation=0.050                                                                                 |          |
| --title          |                            | Title of proposal                                                                                                                           | Yes      |
| --type           |                            | ProposalType of proposal,eg:PlainText/ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage/AddAsset                                                                   | Yes      |
| --version           |            0                | the version of the new protocol                                                                            |       |
| --software           |           " "                 | the software of the new protocol                                                                         |       |
| --switch-height           |       0                     | the switch height of the new protocol                                                         |       |
| --threshold | "0.8"   |  the upgrade signal threshold of the software upgrade                                                   |               |
| --asset-family |  | the asset family, valid values can be fungible and non-fungible | |
| --asset-symbol |  | the asset symbol. Once created, it cannot be modified | |
| --asset-name |  | the asset name | |
| --asset-decimal |  | the asset decimal. The maximum value is 18 | |
| --asset-symbol-min-alias |  | the asset symbol minimum alias | |
| --asset-initial-supply |  | the initial supply token of asset | |
| --asset-max-supply |  | the max supply token of asset | |
| asset-mintable |  | whether the asset can be minted, default false | |

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

### Submit a 'AddAsset' type proposal

```shell
iriscli gov submit-proposal --chain-id=irishub-test --from=node0 --fee=4iris --type=AddAsset --description=test --title=test-proposal --deposit=3000iris --commit --home=$iris_root_path --asset-decimal=18 --asset-family=fungible --asset-initial-supply=100000000 --asset-max-supply=2000000000 --asset-mintable=true --asset-name=IETH --asset-symbol=ETH --asset-symbol-min-alias=eth-atto

```

###  How to query proposal

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
