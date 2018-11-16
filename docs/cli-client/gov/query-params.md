# iriscli gov query-params

## Description

Query parameter proposal's config

## Usage

```
iriscli gov query-params [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 |          |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --key           |                            | [string] key name of parameter                                                                                                                       |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --module        |                            | [string] module name                                                                                                                                 |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

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