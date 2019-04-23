****# iriscli bank

## Description

Bank allows you to manage assets in your local account 

## Usage:

```
 iriscli bank <command>
```


## Available Commands

| Name      | Description                         |
| --------- | ----------------------------------- |
| [coin-type](coin-type.md) | Query coin type                     |
| [token-stats](token-stats.md) | Query token stats     |
| [account](account.md)   | Query account balance               |
| [send](send.md)      | Create/sign/broadcast a send tx           |
| [burn](burn.md)     | Burn tokens |
| [sign](sign.md)     | Sign transactions generated offline |
| [broadcast](broadcast.md)|Broadcast a signed transaction to the network|

## Flags

| Name,shorthand | Default | Description   | Required |
| -------------- | ------- | ------------- | -------- |
| -h, --help     |         | Help for bank |          |

## Global Flags

| Name,shorthand        | Default        | Description                                 | Required |
| --------------------- | -------------- | ------------------------------------------- | -------- |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) |          |
| --home string         | /root/.iriscli | Directory for config and data               |          |
| -o, --output string   | text           | Output format (text \|json)                 |          |
| --trace               |                | Print out full stack trace on errors        |          |
