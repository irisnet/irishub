# iriscli upgrade query-switch

## Description
Query the switch information to know if someone have send the switch message for the certain software upgrade proposal.

## Usage

```
iriscli upgrade query-switch --proposal-id <proposalID> --voter <voter address>
```

## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --proposal-id      |                            | proposalID of upgrade swtich being queried                              | Yes      |
| --voter     |                            | Address sign the switch msg                              | Yes      |
| --chain-id      |                            | [string] Chain ID of tendermint node                              |            |
| --height        | most recent provable block | block height to query                                             |          |
| --help, -h      |                            | help for query                                                    |          |
| --indent        |                            | Add indent to JSON response                                       |          |
| --ledger        |                            | Use a connected Ledger device                                     |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node    | true                       | Don't verify proofs for responses                                 |          |

## Example

Query if the address `faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe` send the switch message for the software upgrade proposal whose ID is 5.

```
iriscli upgrade query-switch --proposal-id=5 --voter=faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe
```
