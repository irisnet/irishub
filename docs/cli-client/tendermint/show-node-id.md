# iris tendermint show-node-id

## Description

Show the id of your node, this id will be used for making connection between peers.

## Usage

```
iris tendermint show-node-id [flags]
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

### Show Node ID of Your Node

```shell
iris tendermint show-node-id --home={iris-home}
```

The sample output could be:
```$xslt
b18d3d1990c886555241f91331f9c00fe69421aa
```

