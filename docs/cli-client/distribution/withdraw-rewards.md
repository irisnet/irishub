# iriscli distribution withdraw-rewards

## Description

Withdraw rewards

## Usage

```
iriscli distribution withdraw-rewards [flags]
```

Print all supported options:

```shell
iriscli distribution withdraw-rewards --help
```

## Unique Flags

| Name, shorthand       | type   | Required | Default  | Description                                                         |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string | false    | ""       | only withdraw from this validator address (in bech) |
| --is-validator        | bool   | false    | false    | Also withdraw validator's commission |

Keep in mind, don't specify the above options both.

## Examples

1. Only withdraw the delegation reward from a given validator
    ```
    iriscli distribution withdraw-rewards --only-from-validator fva134mhjjyyc7mehvaay0f3d4hj8qx3ee3w3eq5nq --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
2. Withdraw all delegation reward of a delegator
    ```
    iriscli distribution withdraw-rewards --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
3. If the delegator is a onwer of a validator, withdraw all delegation reward and validator reward:
    ```
    iriscli distribution withdraw-rewards --is-validator=true --from mykey --fee=0.004iris --chain-id=irishub-test
    ```