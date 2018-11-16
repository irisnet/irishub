# iriscli gov pull-params

## Description

Generate param.json file

## Usage

```
iriscli gov pull-params [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 |          |
| --height        |                            | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for pull-params                                                                                                                                 |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --path          | $HOME/.iris                | [string] Directory of iris home                                                                                                                      |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Pull params

```shell
iriscli gov pull-params
```

Then you'll receive a message as described below:

```txt
Save the parameter config file in  /Users/trevorfu/.iris/config/params.json
```

If you open the params.json in the --path/config directory, you can see it's json format content.

```txt
{
  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
:  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
      "voting_period": "172800000000000"
:  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
      "voting_period": "172800000000000"
    },
    "Gov/govTallyingProcedure": {
      "threshold": "0.5000000000",
      "veto": "0.3340000000",
      "participation": "0.6670000000"
    }
  }
}
```