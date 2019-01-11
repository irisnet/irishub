# iriscli distribution withdraw-rewards

## Description

Withdraw rewards

## Usage

```
iriscli distribution withdraw-rewards [flags]
```

Print help messages:

```
iriscli distribution withdraw-rewards --help
```

## Unique Flags

| Name, shorthand       | type   | Required | Default  | Description                                                         |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string | false    | ""       | only withdraw from this validator address (in bech) |
| --is-validator        | bool   | false    | false    | Also withdraw validator's commission |
| --commit         | String | False     | True                  |wait for transaction commit accomplishment, if true, --async will be ignored|

Keep in mind, don't specify the above options both.

## Examples

1. Only withdraw a delegation rewards from a given validator
    ```
    iriscli distribution withdraw-rewards --only-from-validator <validator address> --from <key name> --fee=0.004iris --chain-id=<chain-id>
    ```
2. Withdraw all delegation rewards of a delegator
    ```
    iriscli distribution withdraw-rewards --from <key name> --fee=0.004iris --chain-id=<chain-id>
    ```
3. If the delegator is a onwer of a validator, withdraw all delegation rewards and validator commmission rewards:
    ```
    iriscli distribution withdraw-rewards --is-validator=true --from <key name> --fee=0.004iris --chain-id=<chain-id>
    ```