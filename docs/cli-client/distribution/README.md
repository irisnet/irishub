# iriscli distribution 

## introduction

This document describe how to use the the command line interfaces of distribution module.

## Usage

```shell
iriscli distribution [subcommand] [flags]
```

Print all supported subcommands and flags:

```shell
iriscli distribution --help
```

## Available Subommands

| Name                            | Description                                                   |
| --------------------------------| --------------------------------------------------------------|
| [delegation-distr-info](delegation-distr-info.md) | Query delegation distribution information |
| [delegator-distr-info](delegator-distr-info.md) | Query delegator distribution information |
| [validator-distr-info](validator-distr-info.md) | Query validator distribution information |
| [withdraw-address](withdraw-address.md) | Query withdraw address |
| [set-withdraw-address](set-withdraw-address.md)  | change the default withdraw address for rewards associated with an address |
| [withdraw-rewards](withdraw-rewards.md) | withdraw rewards for either: all-delegations, a delegation, or a validator |
