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
| --module        |                            | Module name                                                                                                                                 |          |

## Examples
 
### Query params by module

Get all the params of stake module.

```shell
iriscli gov query-params --module=stake
```

```txt
Stake Params:
  Unbonding Time:         504h0m0s
  Max Validators:         100
```


### Modules with all configurable parameters

```shell
iriscli gov query-params --module=auth
iriscli gov query-params --module=mint
iriscli gov query-params --module=stake
iriscli gov query-params --module=slashing
iriscli gov query-params --module=distr
iriscli gov query-params --module=service
```
