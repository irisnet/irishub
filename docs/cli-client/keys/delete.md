# iriscli keys delete

## Description

Delete the given key

## Usage

```
iriscli keys delete <name> [flags]
```

## Flags

| Name, shorthand | Default   | Description                                                  | Required |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | help for add                                                 |          |

## Examples

### Delete a given key

```shell
iriscli keys delete MyKey
```

You'll receive a danger warning and be asked to enter a password for your key.

```txt
DANGER - enter password to permanently delete key:
```

After you enter the correct password, you're done with deleting your key.

```txt
Password deleted forever (uh oh!)
```