# iriscli stake create-validator

## Description

Create new validator initialized with a self-delegation on it

## Usage

```
iriscli stake create-validator [flags]
```

Print help messages:
```shell
iriscli stake create-validator --help
```

## Unique Flags

| Name, shorthand              | type   | Required | Default  | Description                                                         |
| ---------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-delegator          | string | true     | ""       | Bech address of the delegator |
| --amount                     | string | true     | ""       | Amount of coins to bond |
| --commission-max-change-rate | float  | true     | 0.0      | The maximum commission change rate percentage (per day)|
| --commission-max-rate        | float  | true     | 0.0      | The maximum commission rate percentage |
| --commission-rate            | float  | true     | 0.0      | The initial commission rate percentage |
| --details                    | string | false    | ""       | Optional details |
| --genesis-format             | bool   | false    | false    | Export the transaction in gen-tx format; it implies --generate-only |
| --identity                   | string | false    | ""       | Optional identity signature (ex. UPort or Keybase) |
| --ip                         | string | false    | ""       | Node's public IP. It takes effect only when used in combination with |
| --moniker                    | string | true     | ""       | Validator name |
| --pubkey                     | string | true     | ""       | Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |
| --website                    | string | false    | ""       | Optional website |

## Examples

### Create new validator

```shell
iriscli stake create-validator --chain-id=<chain-id> --from=<key name> --fee=0.004iris --pubkey=<Validator PubKey> --commission-max-change-rate=0.01 --commission-max-rate=0.2 --commission-rate=0.1 --amount=100iris --moniker=<validator name>
```

After that, you're done with creating a new validator.

