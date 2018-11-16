# iriscli keys add

## Description

Create a new key, or import from seed

## Usage

```
iriscli keys add <name> [flags]
```

## Flags

| Name, shorthand | Default   | Description                                                       | Required |
| --------------- | --------- | ----------------------------------------------------------------- | -------- |
| --account       |           | [uint32] Account number for HD derivation                         |          |
| --dry-run       |           | Perform action, but don't add key to local keystore               |          |
| --help, -h      |           | help for add                                                      |          |
| --index         |           | [uint32] Index number for HD derivation                           |          |
| --ledger        |           | Store a local reference to a private key on a Ledger device       |          |
| --no-backup     |           | Don't print out seed phrase (if others are watching the terminal) |          |
| --recover       |           | Provide seed phrase to recover existing key instead of creating   |          |
| --type, -t      | secp256k1 | [string] Type of private key (secp256k\|ed25519)                  |          |

## Examples

### Create a new key

```shell
iriscli keys add MyKey
```

You'll be asked to enter a password for your key, note: password must be at least 8 characters.

```txt
Enter a passphrase for your key:
Repeat the passphrase:
```

After that, you're done with creating a new key, but remember to backup your seed phrase, it is the only way to recover your account if you ever forget your password or lose your key.

```txt
NAME:	TYPE:	ADDRESS:						PUBKEY:
MyKey	local	faa1mmsm487rqkgktl2qgrjad0z3yaf9n8t5pkp33m	fap1addwnpepq2g0u7cnxp5ew0yhqep8j4rth5ugq8ky7gjmunk8tkpze95ss23ak4svkjq
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

oval green shrug term already arena pilot spirit jump gain useful symbol hover grid item concert kiss zero bleak farm capable peanut snack basket
```

The 24 words above is a seed phrase just for example, **DO NOT** use it in production.

### Recover an existing key

If you forget your password or lose your key, or you wanna use your key in another place, you can recover your key by your seed phrase.

```txt
iriscli keys add MyKey --recover
```

You'll be asked to enter a new password for your key, and enter the seed phrase. Then you get your key back.

```txt
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

