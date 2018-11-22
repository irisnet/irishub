# iriscli bank

## Description

Bank allows you to manage assets in your local account 

## Usage:

```
 iriscli bank [command]
```

 

## Available Commands

| Name      | Description                         |
| --------- | ----------------------------------- |
| coin-type | Query coin type                     |
| account   | Query account balance               |
| send      | Create and sign a send tx           |
| sign      | Sign transactions generated offline |

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
