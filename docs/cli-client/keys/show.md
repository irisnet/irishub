# iriscli keys show

## Description

Return public details of one local key.

## Usage

```
iriscli keys show <name> [flags]
```

## Flags

| Name, shorthand      | Default           | Description                                                    | Required |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --address            |                   | output the address only (overrides --output)                   |          |
| --bech               | acc               | [string] The Bech32 prefix encoding for a key (acc|val|cons)   |          |
| --help, -h           |                   | help for show                                                  |          |
| --multisig-threshold | 1                 | [uint] K out of N required signatures                          |          |
| --pubkey             |                   | output the public key only (overrides --output)                |          |

## Examples

### Show a given key

```shell
iriscli keys show MyKey
```

You'll get the local public keys with 'address' and 'pubkey' element of a given key.

```txt
NAME:	TYPE:	ADDRESS:						            PUBKEY:
MyKey	local	faa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zakl7lxn5	fap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskww57lw
```
