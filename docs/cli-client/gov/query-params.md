# iriscli gov query-params

## Description

Query parameter proposal's config

## Usage

```
iriscli gov query-params <flags>
```

Print help messages:

```
iriscli gov query-params --help
```
## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --key           |                            | Key name of parameter                                                                                                                       |          |
| --module        |                            | Module name                                                                                                                                 |          |

## Examples
 
### Query params by module

Get all the params of stake module.

```shell
iriscli gov query-params --module=stake
```

```txt
 stake/MaxValidators=100
 stake/UnbondingTime=504h0m0s
```

### Query params by key

Get the details of the parameter specified in the stake module.

```shell
iriscli gov query-params --key=stake/MaxValidators
```

```txt
 stake/MaxValidators=100
```

Note: --module and --key cannot be both empty.

### Modules with all configurable parameters

```shell
iriscli gov query-params --module=auth
iriscli gov query-params --module=mint
iriscli gov query-params --module=stake
iriscli gov query-params --module=slashing
iriscli gov query-params --module=distr
iriscli gov query-params --module=service
```
