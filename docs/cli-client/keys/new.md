# iriscli keys new

## Description

Derive a new private key using an interactive command that will prompt you for each input.
Optionally specify a bip39 mnemonic, a bip39 passphrase to further secure the mnemonic,
and a bip32 HD path to derive a specific account. The key will be stored under the given name
and encrypted with the given password. The only input that is required is the encryption password.

## Usage

```
iriscli keys new <name> [flags]
```

## Flags

| Name, shorthand | Default           | Description                                                     | Required |
| --------------- | ----------------- | --------------------------------------------------------------- | -------- |
| --bip44-path    | 44'/118'/0'/0/0   | BIP44 path from which to derive a private key                   |          |
| --default       |                   | Skip the prompts and just use the default values for everything |          |
| --help, -h      |                   | help for new                                                    |          |
| --ledger        |                   | Store a local reference to a private key on a Ledger device     |          |

## Examples

### Create a new key by the specified method

```shell
iriscli keys new MyKey
```

You'll be asked to enter your bip44 path, default is 44'/118'/0'/0/0.

```txt
> -------------------------------------
> Enter your bip44 path. Default is 44'/118'/0'/0/0
```

Then you'll be asked to enter your bip39 mnemonic, or hit enter to generate one.

```txt
> Enter your bip39 mnemonic, or hit enter to generate one.
```

You could hit enter to generate bip39 mnemonic, then a new hint will be show to ask you to enter bip39 passphrase.

```txt
> -------------------------------------
> Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed
> Most users should just hit enter to use the default, ""
```

Also you could hit enter to skip it, then you'll receive a hint to enter a password.

```txt
> -------------------------------------
> Enter a passphrase to encrypt your key to disk:
> Repeat the passphrase:
```

After that, you're done with creating a new key.
