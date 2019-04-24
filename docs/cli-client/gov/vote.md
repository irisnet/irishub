# iriscli gov vote

## Description

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

## Usage

```
iriscli gov vote <flags>
```

Print help messages:

```
iriscli gov vote --help
```
## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --option         |                            | Vote option {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --proposal-id    |                            | ProposalID of proposal voting on                                                                                                            | Yes      |

## Examples

### Vote for proposal


```shell
iriscli gov vote --chain-id=<chain-id> --proposal-id=<proposal-id> --option=Yes --from=<key_name> --fee=0.3iris
```

Only validators can vote for proposals which enter voting period.

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
     "voter": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul"
   }
 }) 
```

### How to query vote

[query-vote](query-vote.md)

[query-votes](query-votes.md)

[query-tally](query-tally.md)