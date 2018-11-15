# iriscli keys update

## Description

Change the password used to protect private key

## Usage

```
iriscli keys update <name> [flags]
```

## Flags

| Name, shorthand | Default   | Description                                                  | Required |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | help for add                                                 |          |

## Examples

### Change the password of a given key

```shell
iriscli keys update MyKey
```

You'll be asked to enter the current password for your key.

```txt
Enter the current passphrase:
```

After this, you'll be asked to enter a new password.

```txt
Enter the new passphrase:
Repeat the new passphrase:
```

After that, you're done with changing your password.

```txt
Password successfully updated!
```