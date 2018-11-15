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

You'll be asked to enter a password for your key.

```txt
DANGER - enter password to permanently delete key:
```

After that, you're done with deleting your key.

```txt
Password deleted forever (uh oh!)
```