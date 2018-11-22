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
iriscli gov query-params --module=gov
```

You'll get all the keys of gov module.

```txt
[
 "Gov/govDepositProcedure",
 "Gov/govTallyingProcedure",
 "Gov/govVotingProcedure"
]
```

### Query params by key

```shell
iriscli gov query-params --key=Gov/govDepositProcedure
```

You'll get all the details of the key specified in the gov module.

```txt
{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"1000000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":""}
```

Note: --module and --key cannot be both empty.
