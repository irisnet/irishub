# iriscli service

## Description
Service allows you to define、bind、invocate a service on chain.

## Usage

```shell
iriscli service [command]
```

## Available Commands

| Name                                  | Description                               |
| ------------------------------------  | ----------------------------------------- |
| [define](define.md)                   | Create a new service definition           |
| [definition](definition.md)           | Query service definition                  |
| [bind](bind.md)                       | Create a new service binding              |
| [binding](binding.md)                 | Query service binding                     |
| [bindings](bindings.md)               | Query service bindings                    |
| [update-binding](update-binding.md)   | Update a service binding                  |
| [disable](disable.md)                 | Disable a available service binding       |
| [enable](enable.md)                   | Enable an unavailable service binding     |
| [refund-deposit](refund-deposit.md)   | Refund all deposit from a service binding |

## Flags

| Name, shorthand | Default | Description      | Required |
| --------------- | ------- | ---------------- | -------- |
| --help, -h      |         | help for service |          |