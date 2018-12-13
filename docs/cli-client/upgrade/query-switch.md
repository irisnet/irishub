# iriscli upgrade query-switch

## Description
Query the switch information to know if someone have send the switch message for the certain software upgrade proposal.

## Usage

```
iriscli upgrade query-switch --proposal-id <proposalID> --voter <voter address>
```

Print help messages:

```
iriscli upgrade query-switch --help
```
## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --proposal-id     |                            | proposalID of upgrade swtich being queried                              | Yes      |
| --voter     |                            | Address sign the switch msg                              | Yes      |

## Example

Query if the address `faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe` send the switch message for the software upgrade proposal whose ID is 5.

```
iriscli upgrade query-switch --proposal-id=5 --voter=faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe
```
