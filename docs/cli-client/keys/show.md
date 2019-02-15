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

### Show Validator Operator Address

If an address has bonded to be a validator operator, then you could use `iriscli keys show` to get the operator's 
address:

```$xslt
iriscli keys show alice --bech val
```

Then you could see the following:
```$xslt
NAME: TYPE: ADDRESS: PUBKEY:
alice local fva12nda6xwpmp000jghyneazh4kkgl2tnzy73dmzy fvp1addwnpepqfw52vyzt9xgshxmw7vgpfqrey30668g36f9z837kj9dy68kn2wxqh3zvz9
```

The result could be use for `--address-validator` in [create a delegation](../stake/delegate.md)