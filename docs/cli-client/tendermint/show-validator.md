# iris tendermint show-validator

## Description

Show the bech32 encoded public key of your validator, with hrp of `icp`, which derives from private_key.json in <path_to_your_home>/config

## Usage

```
iris tendermint show-validator <flags>
```

## Flags

| Name, shorthand      | Default           | Description                                                    | Required |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --json            |                   | output in json format                  |          |
| --help, -h           |                   | help for show                                                  |          |

## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Show Public Key of Your Node

```
iris tendermint show-validator --home=<path_to_your_home>
```

The sample output could be:
```
icp1zcjduepqzuz420weqehs3mq0qny54umfk5r78yup6twtdt7mxafrprms5zqsjeuxvx
```

The output is encoded in Bech32, to read more about this encoding method, read [this](../../features/basic-concepts/bech32-prefix.md) 