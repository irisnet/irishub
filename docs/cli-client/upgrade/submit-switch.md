# iriscli upgrade submit-switch

## Description

Submit a switch msg for a upgrade propsal after installing the new software and broadcast to the whole network.

## Usage

```
iriscli upgrade submit-switch [flags]
```

Print help messages:

```
iriscli upgrade submit-switch --help
```
## Flags

| Name, shorthand  | Default   | Description                                                  | Required |
| ---------------  | --------- | ------------------------------------------------------------ | -------- |
| --proposal-id    |           | proposalID of upgrade proposal                               | Yes      |
| --title          |           | title of switch                                              |          |

## Examples

Send a switch message for the software upgrade proposal whose `proposalID` is 5. 

```
iriscli upgrade submit-switch --chain-id=IRISnet --from=x --fee=0.004iris --proposal-id 5 --title="Run new verison"
```
