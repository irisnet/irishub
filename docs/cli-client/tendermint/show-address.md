# iris tendermint show-address

## Description

Shows this node's tendermint validator address

## Usage

```
iris tendermint show-address [flags]
```

## Flags

| Name, shorthand      | Default           | Description                                                    | Required |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --help, -h           |                   | help for show                                                  |          |

## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Show Address of Your Node

```shell
iris tendermint show-address --home={iris-home}
```

The sample output could be:
```$xslt
fca1jhax359kz6hm70swxpxumpgkaglr7yq80qhwvv
```

The output is encoded in Bech32, to read more about this encoding method, read [this](../../features/basic-concepts/bech32-prefix.md)
