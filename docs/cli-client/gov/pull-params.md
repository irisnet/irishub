# iriscli gov pull-params

## Description

Generate param.json file

## Usage

```
iriscli gov pull-params [flags]
```

Print help messages:

```
iriscli gov pull-params --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --path          | $HOME/.iriscli                | [string] Directory of iris home                                                                                                                      |          |

## Examples

### Pull params

```shell
iriscli gov pull-params
```

Then you'll receive a message as described below:

```txt
Save the parameter config file in  /Users/trevorfu/.iriscli/params.json
```

If you open the params.json in the --path directory, you can see it's json format content.

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
