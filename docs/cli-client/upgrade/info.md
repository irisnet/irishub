# iriscli upgrade info

## Description

Query the information of software version and upgrade module.

## Usage

```
iriscli upgrade info
```

## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                              |            |
| --height        | most recent provable block | block height to query                                             |          |
| --help, -h      |                            | help for query                                                    |          |
| --indent        |                            | Add indent to JSON response                                       |          |
| --ledger        |                            | Use a connected Ledger device                                     |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node    | true                       | Don't verify proofs for responses                                 |          |

## Example

Query the current version information. 

```
iriscli upgrade info 
```

Then it will show 

```
{
"current_proposal_id": "0",
"current_proposal_accept_height": "-1",
"version": {
"Id": "0",
"ProposalID": "0",
"Start": "0",
"ModuleList": [
{
"Start": "0",
"End": "9223372036854775807",
"Handler": "bank",
"Store": [
"acc"
]
},
{
"Start": "0",
"End": "9223372036854775807",
"Handler": "stake",
"Store": [
"stake",
"acc",
"mint",
"distr"
]
},
.......
]
}
}
```
