# iriscli gov

This module provides the basic functionalities for [Governance](../features/governance.md).

## Available Commands

| Name                                            | Description                                                     |
| ----------------------------------------------- | --------------------------------------------------------------- |
| [query-proposal](#iriscli-gov-query-proposal)   | Query details of a single proposal                              |
| [query-proposals](#iriscli-gov-query-proposals) | Query proposals by conditions                                   |
| [query-vote](#iriscli-gov-query-vote)           | Query vote                                                      |
| [query-votes](#iriscli-gov-query-votes)         | Query votes on a proposal                                       |
| [query-deposit](#iriscli-gov-query-deposit)     | Query details of a deposit                                      |
| [query-deposits](#iriscli-gov-query-deposits)   | Query deposits on a proposal                                    |
| [query-tally](#iriscli-gov-query-tally)         | Query the statistics of a proposal                              |
| [submit-proposal](#iriscli-gov-submit-proposal) | Submit a proposal along with an initial deposit                 |
| [deposit](#iriscli-gov-deposit)                 | Deposit tokens for an active proposal                           |
| [vote](#iriscli-gov-vote)                       | Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain |

## iriscli gov query-proposal

Query details of a proposal

```bash
iriscli gov query-proposal <flags>
```

**Flags:**

| Name, shorthand | Type | Required | Default | Description            |
| --------------- | ---- | -------- | ------- | ---------------------- |
| --proposal-id   | uint | Yes      |         | Identity of a proposal |

### Query a proposal

```bash
iriscli gov query-proposal --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-proposals

Query proposals by conditions

```bash
iriscli gov query-proposals <flags>
```

**Flags:**

| Name, shorthand | Type    | Required | Default | Description                                                         |
| --------------- | ------- | -------- | ------- | ------------------------------------------------------------------- |
| --depositor     | Address |          |         | Filter proposals by depositor address                               |
| --limit         | uint    |          |         | Limit to the latest [number] of proposals. Default to all proposals |
| --status        | string  |          |         | Filter proposals by status (passed / rejected)                      |
| --voter         | Address |          |         | Filter proposals by voter address                                   |

### Query all proposals

```bash
iriscli gov query-proposals --chain-id=irishub
```

### Query proposals by conditions

```bash
iriscli gov query-proposals --chain-id=irishub --limit=3 --status=passed --depositor=<iaa...>
```

## iriscli gov query-vote

Query a vote

```bash
iriscli gov query-vote <flags>
```

**Flags:**

| Name, shorthand | Type | Required | Default | Description            |
| --------------- | ---- | -------- | ------- | ---------------------- |
| --proposal-id   |  uint    | Yes      |         | Identity of a proposal |
| --voter         |  Address    | Yes      |         | Bech32 voter address   |

### Query a vote

```bash
iriscli gov query-vote --chain-id=irishub --proposal-id=<proposal-id> --voter=<iaa...>
```

## iriscli gov query-votes

Query all votes of a proposal

```bash
iriscli gov query-votes <flags>
```

**Flags:**

| Name, shorthand | Type | Required | Default | Description            |
| --------------- | ---- | -------- | ------- | ---------------------- |
| --proposal-id   |  uint    | Yes      |         | Identity of a proposal |

### Query all votes of a proposal

```bash
iriscli gov query-votes --chain-id=irishub --proposal-id=<proposal-id>

```

## iriscli gov query-deposit

Query a deposit of a proposal

```bash
iriscli gov query-deposit <flags>

```

**Flags:**

| Name, shorthand | Type | Required | Default | Description              |
| --------------- | ---- | -------- | ------- | ------------------------ |
| --proposal-id   | uint     | Yes      |         | Identity of a proposal   |
| --depositor     |  Address    | Yes      |         | Bech32 depositor address |

### Query a deposit of a proposal

```bash
iriscli gov query-deposit --chain-id=irishub --proposal-id=<proposal-id> --depositor=<iaa...>

```

## iriscli gov query-deposits

Query all deposits of a proposal

```bash
iriscli gov query-deposits <flags>

```

**Flags:**

| Name, shorthand | Type | Required | Default | Description            |
| --------------- | ---- | -------- | ------- | ---------------------- |
| --proposal-id   |   uint   | Yes      |         | Identity of a proposal |

### Query all deposits of a proposal

```bash
iriscli gov query-deposits --chain-id=irishub --proposal-id=<proposal-id>

```

## iriscli gov query-tally

Query the statistics of a proposal

```bash
iriscli gov query-tally <flags>

```

**Flags:**

| Name, shorthand | Type | Required | Default | Description            |
| --------------- | ---- | -------- | ------- | ---------------------- |
| --proposal-id   |  uint    | Yes      |         | Identity of a proposal |

### Query the statistics of a proposal

```bash
iriscli gov query-tally --chain-id=irishub --proposal-id=<proposal-id>

```

## iriscli gov submit-proposal

Submit a proposal along with an initial deposit

```bash
iriscli gov submit-proposal <flags>

```

**Flags:**

| Name, shorthand          | Type   | Required | Default | Description                                                                                                    |
| ------------------------ | ------ | -------- | ------- | -------------------------------------------------------------------------------------------------------------- |
| --deposit                | Coin   | Yes      |         | Initial deposit of the proposal(at least  30% of minDeposit)                                                   |
| --description            | string | Yes      |         | Description of the proposal                                                                                    |
| --param                  | string |          |         | On-chain Parameter to be changed, eg. mint/Inflation=0.050                                                     |
| --title                  | string | Yes      |         | Title of the proposal                                                                                          |
| --type                   | string | Yes      |         | ProposalType of the proposal(PlainText/Parameter/SoftwareUpgrade/SoftwareHalt/CommunityTaxUsage/TokenAddition) |
| --version                | uint   |          | 0       | The version of the new protocol                                                                                |
| --software               | string |          |         | The software of the new protocol                                                                               |
| --switch-height          | uint   |          | 0       | The switch height of the new protocol                                                                          |
| --threshold              | string |          | "0.8"   | The upgrade signal threshold of the software upgrade                                                           |
| --token-canonical-symbol | string |          |         | The source symbol of a external token                                                                          |
| --token-symbol           | string |          |         | The token symbol. Once created, it cannot be modified                                                          |
| --token-name             | string |          |         | The token name                                                                                                 |
| --token-decimal          | uint   |          |         | The token decimal. The maximum value is 18                                                                     |
| --token-min-unit-alias   | string |          |         | The token symbol minimum alias                                                                                 |
| --token-initial-supply   | uint64 |          |         | The initial supply token of token                                                                              |

:::tip
The proposer must deposit at least 30% of the [MinDeposit](../features/governance.md#proposal-level) to submit a proposal.
:::

### Submit a Parameter Change Proposal

:::tip
[What parameters can be changed online?](../concepts/gov-params.md)
:::

**Unique Required Params:** `--param`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=0.3iris --deposit=2000iris --type=Parameter --param='mint/Inflation=0.050'

```

### Submit a Software Upgrade Proposal

**Unique Required Params:** `--software`, `--version`, `--switch-height`, `--threshold`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=0.3iris --deposit=2000iris --type=SoftwareUpgrade --software=https://github.com/irisnet/irishub/tree/v0.15.1 --version=2 --switch-height=8000 --threshold=0.8

```

### Submit a Token Addition Proposal

**Unique Params:**

- Required: `--token-symbol`, `--token-canonical-symbol`, `--token-name`
- Optional: `--token-decimal`, `--token-min-unit-alias`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=1iris --deposit=2000iris --type=TokenAddition --token-symbol=btc --token-canonical-symbol=btc --token-name=Bitcoin --token-decimal=18 --token-min-unit-alias=satoshi

```

## iriscli gov deposit

Deposit tokens for an active proposal

```bash
iriscli gov deposit <flags>

```

**Flags:**

| Name, shorthand | Type | Required | Default | Description             |
| --------------- | ---- | -------- | ------- | ----------------------- |
| --deposit       | Coin | Yes      |         | Deposit of the proposal |
| --proposal-id   | uint | Yes      |         | Identity of a proposal  |

### Deposit for an active proposal

When the total deposit amount exceeds the [MinDeposit](../features/governance.md#proposal-level), the proposal will enter the voting procedure.

```bash
iriscli gov deposit --chain-id=irishub --proposal-id=<proposal-id> --deposit=50iris --from=<key-name> --fee=0.3iris
```

## iriscli gov vote

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

:::tip
[No VS NoWithVeto](../features/governance.md#burning-mechanism)

Only validators and delegators can vote for proposals in the voting period.
:::

```bash
iriscli gov vote <flags>
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                            |
| --------------- | ------ | -------- | ------- | -------------------------------------- |
| --option        | string | Yes      |         | Vote option: Yes/No/NoWithVeto/Abstain |
| --proposal-id   | uint   | Yes      |         | Identity of a proposal                 |

### Vote for an active proposal

```bash
iriscli gov vote --chain-id=irishub --proposal-id=<proposal-id> --option=Yes --from=<key-name> --fee=0.3iris
```
