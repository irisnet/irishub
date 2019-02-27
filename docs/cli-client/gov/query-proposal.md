# iriscli gov query-proposal

## Description

Query details of a single proposal

## Usage

```
iriscli gov query-proposal [flags]
```

Print help messages:

```
iriscli gov query-proposal --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query proposal

```shell
iriscli gov query-proposal --chain-id=<chain-id> --proposal-id=1
```

You could query the details of a specific proposal.

```txt
{
  "proposal_id": "1",
  "title": "test proposal",
  "description": "a new text proposal",
  "proposal_type": "Text",
  "proposal_status": "DepositPeriod",
  "tally_result": {
    "yes": "0.0000000000",
    "abstain": "0.0000000000",
    "no": "0.0000000000",
    "no_with_veto": "0.0000000000"
  },
  "submit_time": "2018-11-14T09:10:19.365363Z",
  "deposit_end_time": "2018-11-16T09:10:19.365363Z",
  "total_deposit": [
    {
      "denom": "iris-atto",
      "amount": "49000000000000000050"
    }
  ],
  "voting_start_time": "0001-01-01T00:00:00Z",
  "voting_end_time": "0001-01-01T00:00:00Z",
  "param": {
    "key": "",
    "value": "",
    "op": ""
  }
}
```
