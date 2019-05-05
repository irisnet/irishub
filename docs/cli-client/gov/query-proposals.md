# iriscli gov query-proposals

## Description

Query proposals with optional filters

## Usage

```
iriscli gov query-proposals <flags>
```


Print help messages:

```
iriscli gov query-proposals --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --depositor     |                            | Filter by proposals deposited on by depositor                                                                                    |  false     |
| --limit         |                            | Limit to latest [number] proposals. Default to all proposals                                                                    |    false      |
| --status        |                            | filter proposals by proposal status                                                                                                        |    false      |
| --voter         |                            | Filter by proposals voted on by voted                                                                                            |     false     |

## Examples

### Query proposals

```shell
iriscli gov query-proposals --chain-id=<chain-id>
```

You could query all the proposals by default.

```txt
  1 - test proposal
  2 - new proposal
```

Also you can query proposal by filters, such as:

```shell
gov query-proposals --chain-id=<chain-id> --depositor=iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs
```

Finally, here shows the proposal who's depositor address is iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs.

```txt
  2 - new proposal
```

Query latest 3 proposals
```shell
iriscli gov query-proposals --chain-id=<chain-id> --limit=3
```