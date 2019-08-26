# iriscli rand

## Description

this module allows you to post a random number request to the IRIS Hub and query the random numbers or the pending random number requests

## Usage

```bash
iriscli rand <command>
```

Print all supported subcommands and flags:

```bash
iriscli rand --help
```

## Available Commands

| Name                            | Description                                                      |
| ------------------------------- | ---------------------------------------------------------------- |
| [request-rand](request-rand.md) | Request a random number                                          |
| [query-rand](query-rand.md)     | Query the generated random number by the request id              |
| [query-queue](query-queue.md)   | Query the pending random number requests with an optional height |
