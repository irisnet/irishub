# iriscli asset

## Description

Asset allows you to manage assets on IRIS Hub

## Usage:

```
 iriscli asset <command>
```


## Available Commands

| Name                                                | Description                                     |
| --------------------------------------------------- | ----------------------------------------------- |
| [create-gateway](create-gateway.md)                 | Create a gateway                                |
| [edit-gateway](edit-gateway.md)                     | Edit a gateway                                  |
| [transfer-gateway-owner](transfer-gateway-owner.md) | Transfer the ownership of a gateway             |
| [issue-token](issue-token.md)                       | Issue a new token                               |
| [edit-token](edit-token.md)                         | Edit an existing token                          |
| [transfer-token-owner](transfer-token-owner.md)     | Transfer the ownership of a token               |
| [mint-token](mint-token.md)                         | Mint tokens to a specified address              |
| [query-token](query-token.md)                       | Query details of a token                        |
| [query-tokens](query-tokens.md)                     | Query details of a group of tokens              |
| [query-gateway](query-gateway.md)                   | Query details of a gateway by the given moniker |
| [query-gateways](query-gateways.md)                 | Query all gateways with an optional owner       |
| [query-fee](query-fee.md)                           | Query the asset related fees                    |


## Global Flags

| Name, shorthand       | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding        | hex            | String binary encoding (hex\b64\btc )       |          | string |
| --home                | /root/.iriscli | Directory for config and data               |          | string |
| -o, --output          | text           | Output format (text\json)                   |          | string |
| --trace               |                | Print out full stack trace on errors        |          |        |

## Flags

| Name, shorthand | Default | Description    | Required |
| --------------- | ------- | -------------- | -------- |
| -h, --help      |         | Help for asset |          |
