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
Committed at block 43 (tx hash: 01C4C3B00C6048A12AE2CF2294F63C55A69011381B819C35F11B04C921DB81CC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 2048,
   "codespace": "",
   "tags": {
     "action": "vote",
     "proposal-id": "2",
     "voter": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz"
   }
 })
```
