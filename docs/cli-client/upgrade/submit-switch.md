# iriscli upgrade submit-switch

## Description

Submit a switch msg for a upgrade propsal after installing the new software and broadcast to the whole network.

## Usage

```
iriscli upgrade submit-switch [flags]
```

## Flags

| Name, shorthand  | Default   | Description                                                  | Required |
| ---------------  | --------- | ------------------------------------------------------------ | -------- |
| --proposal-id    |           | proposalID of upgrade proposal                               | Yes      |
| --title          |           | title of switch                                              |          |
| --help, -h       |           | help for submit-switch                                       |          |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                   |          |
| --async          |                            | Broadcast transactions asynchronously                                                       |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                        | Yes      |
| --dry-run        |                            | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it     |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                  | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                             | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                         |          |
| --gas string     | 200000                     | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |          |
| --gas-adjustment | 1                          | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                        |          |
| --indent         |                            | Add indent to JSON response                                                                 |          |
| --json           |                            | return output in json format                                                                |          |
| --ledger         |                            | Use a connected Ledger device                                                               |          |
| --memo           |                            | [string] Memo to send along with transaction                                                |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                           |                    | Yes      |
| --print-response |                            | return tx response (only works with async = false)                                          |          |
| --sequence       |                            | [int] Sequence number to sign the tx                                                        |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                           |          |
## Examples

Send a switch message for the software upgrade proposal whose `proposalID` is 5. 

```
iriscli upgrade submit-switch --chain-id=IRISnet --from=x --fee=0.004iris --proposalID 5 --title="Run new verison"
```
