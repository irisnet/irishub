# iriscli gov query-proposals

## Description

Query proposals with optional filters

## Usage

```
iriscli gov query-proposals [flags]
```


Print help messages:

```
iriscli gov query-proposals --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --depositor     |                            | [string] (optional) Filter by proposals deposited on by depositor                                                                                    |          |
| --limit         |                            | [string] (optional) Limit to latest [number] proposals. Defaults to all proposals                                                                    |          |
| --status        |                            | [string] (optional) filter proposals by proposal status                                                                                                        |          |
| --voter         |                            | [string] (optional) Filter by proposals voted on by voted                                                                                            |          |

## Examples

### Query proposals

```shell
iriscli gov query-proposals --chain-id=test
```

You could query all the proposals by default.

```txt
  1 - test proposal
  2 - new proposal
```

Also you can query proposal by filters, such as:

```shell
gov query-proposals --chain-id=test --depositor=iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

Finally, here shows the proposal who's depositor address is iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd.

```txt
  2 - new proposal
```
