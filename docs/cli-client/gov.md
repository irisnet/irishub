# iriscli gov

This module provides the basic functionalities for [Governance](../features/governance.md).

## Available Commands

| Name                                  | Description                                                     |
| ------------------------------------- | --------------------------------------------------------------- |
| [query-proposal](query-proposal.md)   | Query details of a single proposal                              |
| [query-proposals](query-proposals.md) | Query proposals by conditions                           |
| [query-vote](query-vote.md)           | Query vote                                                      |
| [query-votes](query-votes.md)         | Query votes on a proposal                                       |
| [query-deposit](query-deposit.md)     | Query details of a deposit                                      |
| [query-deposits](query-deposits.md)   | Query deposits on a proposal                                    |
| [query-tally](query-tally.md)         | Query the statistics of a proposal                                |
| [submit-proposal](submit-proposal.md) | Submit a proposal along with an initial deposit                 |
| [deposit](deposit.md)                 | Deposit tokens for an active proposal                              |
| [vote](vote.md)                       | Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain |

## iriscli gov query-proposal

Query details of a proposal

```bash
iriscli gov query-proposal <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --proposal-id   |         | Identity of a proposal               | true     |

### Query a proposal

```bash
iriscli gov query-proposal --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-proposals

Query proposals by conditions

```bash
iriscli gov query-proposals <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                                                  | Required |
| --------------- | ------- | ------------------------------------------------------------ | -------- |
| --depositor     |         | Filter proposals by depositor address                |          |
| --limit         |         | Limit to the latest [number] of proposals. Default to all proposals |          |
| --status        |         | Filter proposals by status (passed / rejected)                          |          |
| --voter         |         | Filter proposals by voter address                |          |

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

**Unique Flags:**

| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --proposal-id   |         | Identity of a proposal | true     |
| --voter         |         | Bech32 voter address                 | true     |

### Query a vote

```bash
iriscli gov query-vote --chain-id=irishub --proposal-id=<proposal-id> --voter=<iaa...>
```

## iriscli gov query-votes

Query all votes of a proposal

```bash
iriscli gov query-votes <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --proposal-id   |         | Identity of a proposal | true      |

### Query all votes of a proposal

```bash
iriscli gov query-votes --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-deposit

Query a deposit of a proposal

```bash
iriscli gov query-deposit <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                 | Required |
| --------------- | ------- | --------------------------- | -------- |
| --proposal-id   |         | Identity of a proposal | true     |
| --depositor     |         | Bech32 depositor address    | true     |

### Query a deposit of a proposal

```bash
iriscli gov query-deposit --chain-id=irishub --proposal-id=<proposal-id> --depositor=<iaa...>
```

## iriscli gov query-deposits

Query all deposits of a proposal

```bash
iriscli gov query-deposits <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --proposal-id   |         | Identity of a proposal | true     |

### Query all deposits of a proposal

```bash
iriscli gov query-deposits --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-tally

Query the statistics of a proposal

```bash
iriscli gov query-tally <flags>
```

**Unique Flags:**
| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --proposal-id   |         | Identity of a proposal | true     |

### Query the statistics of a proposal

```bash
iriscli gov query-tally --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov submit-proposal

Submit a proposal along with an initial deposit

```bash
iriscli gov submit-proposal <flags>
```

**Unique Flags:**

| Name, shorthand          | Default | Description                                                                                                  | Required |
| ------------------------ | ------- | -------------------------------------------------------------------------------------------------------------| -------- |
| --deposit                |         | Initial deposit of the proposal(at least  30% of minDeposit)                                                             | true     |
| --description            |         | Description of the proposal                                                                                      | true     |
| --param                  |         | On-chain Parameter to be changed, eg. mint/Inflation=0.050                                                               |          |
| --title                  |         | Title of the proposal                                                                                            | true     |
| --type                   |         | ProposalType of the proposal (PlainText/Parameter/SoftwareUpgrade/SoftwareHalt/CommunityTaxUsage/TokenAddition) | true     |
| --version                | 0       | The version of the new protocol                                                                              |          |
| --software               |         | The software of the new protocol                                                                             |          |
| --switch-height          | 0       | The switch height of the new protocol                                                                        |          |
| --threshold              | "0.8"   | The upgrade signal threshold of the software upgrade                                                         |          |
| --token-canonical-symbol |         | The source symbol of a external token                                                                        |          | 
| --token-symbol           |         | The token symbol. Once created, it cannot be modified                                                        |          |
| --token-name             |         | The token name                                                                                               |          |
| --token-decimal          |         | The token decimal. The maximum value is 18                                                                   |          |
| --token-min-unit-alias   |         | The token symbol minimum alias                                                                               |          |
| --token-initial-supply   |         | The initial supply token of token                                                                            |          |

:::tip
The proposer must deposit at least 30% of the [MinDeposit](#TODO) to submit a proposal.
:::

### Submit a Parameter Change Proposal

:::tip
[What parameters can be changed online](#TODO)
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

**Unique Flags:**

| Name, shorthand | Default | Description                          | Required |
| --------------- | ------- | ------------------------------------ | -------- |
| --deposit       |         | Deposit of the proposal                  | true     |
| --proposal-id   |         | Identity of a proposal | true     |

### Deposit for an active proposal

When the total deposit amount exceeds the [MinDeposit](#TODO), the proposal will enter the voting procedure.

```bash
iriscli gov deposit --chain-id=irishub --proposal-id=<proposal-id> --deposit=50iris --from=<key-name> --fee=0.3iris
```

## iriscli gov vote

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

:::tip
[No VS NoWithVeto](#TODO)

Only validators and delegators can vote for proposals in the voting period.
:::

```bash
iriscli gov vote <flags>
```

**Unique Flags:**

| Name, shorthand  | Default | Description                                | Required |
| ---------------- | ------- | ------------------------------------------ | -------- |
| --option         |         | Vote option: Yes/No/NoWithVeto/Abstain | true     |
| --proposal-id    |         | Identity of a proposal           | true     |

### Vote for an active proposal

```bash
iriscli gov vote --chain-id=irishub --proposal-id=<proposal-id> --option=Yes --from=<key-name> --fee=0.3iris
```
