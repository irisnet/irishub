# iriscli gov query-params

## Description

Query parameter proposal's config

## Usage

```
iriscli gov query-params [flags]
```

Print help messages:

```
iriscli gov query-params --help
```
## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --key           |                            | [string] Key name of parameter                                                                                                                       |          |
| --module        |                            | [string] Module name                                                                                                                                 |          |

## Examples
 
### Query params by module

```shell
iriscli gov query-params --module=stake
```

You'll get all the params of stake module.

```txt
 stake/MaxValidators=100
 stake/UnbondingTime=504h0m0s
```

### Query params by key

```shell
iriscli gov query-params --key=stake/MaxValidators
```

You'll get the details of the parameter specified in the stake module.

```txt
 stake/MaxValidators=100
```

Note: --module and --key cannot be both empty.