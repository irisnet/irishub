# iris tendermint show-node-id

## Description

This command shows the hex-encoding of address which derives from private key of node_key.json in <path_to_your_home>/config

## Usage

```
iris tendermint show-node-id <flags>
```

## Global Flags

| Name, shorthand       | Default        | Description                             | Required | Type   |
| --------------------- | -------------- | ----------------------------------------| -------- | ------ |
| -e, --encoding string | hex            | String binary encoding (hex\b64\btc )   |          | string |
| --home string         | /root/.iriscli | Directory for config and data           |          | string |
| -o, --output string   | text           | Output format (text\json)               |          | string |
| --trace               |                | Print out full stack trace on errors    |          |        |

## Examples

### Show Node ID of Your Node

```shell
iris tendermint show-node-id --home=<path_to_your_home>
```

The sample output could be:
```
b18d3d1990c886555241f91331f9c00fe69421aa
```

