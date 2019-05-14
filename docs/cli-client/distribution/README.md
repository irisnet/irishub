# iriscli distribution 

## introduction

This document describes how to use the the command line interfaces of distribution module.

## Usage

```
iriscli distribution [subcommand] [flags]
```

Print all supported subcommands and flags:

```
iriscli distribution --help
```

## Available Subommands

| Name                            | Description                                                   |
| --------------------------------| --------------------------------------------------------------|
| [withdraw-address](withdraw-address.md) | Query withdraw address |
| [rewards](rewards.md) | Query all the rewards of validator or delegator |
| [set-withdraw-address](set-withdraw-address.md)  | change withdraw address |
| [withdraw-rewards](withdraw-rewards.md) | withdraw rewards for either: all-delegations, a delegation, or a validator |
| [community-tax](community-tax.md) | Query community tax |