# iriscli gov query-deposit
 ## Description
 Query details of a deposit
 ## Usage
 ```
iriscli gov query-deposit [flags]
```
 ## Flags
| Name, shorthand | Default               | Description                                                                                                                                          | Required |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                       | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                       | [string] bech32 depositer address                                                                                                                    |          |
| --height        |                       | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                       | help for submit-proposal                                                                                                                             |          |
| --indent        |                       | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                       | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                       | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                  | Don't verify proofs for responses                                                                                                                    |          |
 ## Examples
 ### Query deposit
 ```shell
iriscli gov query-deposit --chain-id=test --proposal-id=1 --depositer=faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07
```
 After that, you're done with depositing iris tokens for an activing proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.
 ```txt
{
  "depositer": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
  "proposal_id": "1",
  "amount": [
    {
      "denom": "iris-atto",
      "amount": "30000000000000000000"
    }
  ]
}
```