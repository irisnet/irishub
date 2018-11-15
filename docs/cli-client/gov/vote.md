# iriscli gov vote

## Description

Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain

## Usage

```
iriscli gov vote [flags]
```

## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                                                                            |          |
| --async          |                            | broadcast transactions asynchronously                                                                                                                |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --dry-run        |                            | ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                                                              |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                                                                           | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                                                                                      | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                                                                                  |          |
| --gas            | 200000                     | [string] gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                                                 |          |
| --gas-adjustment | 1                          | [float] adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                                                                                 |          |
| --help, -h       |                            | help for submit-proposal                                                                                                                             |          |
| --indent         |                            | Add indent to JSON response                                                                                                                          |          |
| --json           |                            | return output in json format                                                                                                                         |          |
| --ledger         |                            | Use a connected Ledger device                                                                                                                        |          |
| --memo           |                            | [string] Memo to send along with transaction                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --option         |                            | [string] vote option {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --print-response |                            | return tx response (only works with async = false)                                                                                                   |          |
| --proposal-id    |                            | [string] proposalID of proposal voting on                                                                                                            | Yes      |
| --sequence       |                            | [int] Sequence number to sign the tx                                                                                                                 |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Vote

```shell

```

 After that, you're done with submitting a new proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt

```