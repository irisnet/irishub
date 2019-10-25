# iriscli keys

Keys allows you to manage your local tendermint keystore (wallets) for iris.

## Available Commands

| Name                               | Description                                                                                  |
| ---------------------------------- | -------------------------------------------------------------------------------------------- |
| [add](#iriscli-keys-add)           | Create a new key, or import from seed , or import from a keystore file                       |
| [list](#iriscli-keys-list)         | List all keys                                                                                |
| [show](#iriscli-keys-show)         | Show key info for the given name                                                             |
| [export](#iriscli-keys-export)     | Export keystore to a json file                                                               |
| [delete](#iriscli-keys-delete)     | Delete the given key                                                                         |
| [update](#iriscli-keys-update)     | Change the password used to protect private key                                              |
| [mnemonic](#iriscli-keys-mnemonic) | Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. |
| [new](#iriscli-keys-new)           | Derive a new private key using an interactive command that will prompt you for each input.   |

## iriscli keys add

Create a new key (wallet), or recover from mnemonic/keystore.

```bash
iriscli keys add <key-name> <flags>
```

**Flags:**

| Name, shorthand      | Default   | Description                                                       | Required |
| -------------------- | --------- | ----------------------------------------------------------------- | -------- |
| --account            |           | Account number for HD derivation                                  |          |
| --dry-run            |           | Perform action, but don't add key to local keystore               |          |
| --help, -h           |           | Help for add                                                      |          |
| --index              |           | Index number for HD derivation                                    |          |
| --ledger             |           | Store a local reference to a private key on a Ledger device       |          |
| --no-backup          |           | Don't print out seed phrase (if others are watching the terminal) |          |
| --recover            |           | Provide seed phrase to recover existing key instead of creating   |          |
| --keystore           |           | Recover a key from keystore                                       |          |
| --multisig           |           | Create multisig key                                               |          |
| --multisig-threshold |           | Specify the minimum number of signatures for multisig key         |          |
| --type, -t           | secp256k1 | Type of private key (secp256k\ed25519)                            |          |

### Create a new key

```bash
iriscli keys add MyKey
```

Enter and repeat the password, at least 8 characters, then you will get a new key.

:::warning
**Important**

write the seed phrase in a safe place! It is the only way to recover your account if you ever forget your password.
:::

### Recover an existing key from seed phrase

If you forget your password or lose your key, or you wanna use your key in another place, you can recover your key by your seed phrase.

```bash
iriscli keys add MyKey --recover
```

You'll be asked to enter and repeat the new password for your key, and enter the seed phrase. Then you get your key back.

```bash
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

### Import an existing key from keystore

```bash
iriscli keys add Mykey --recover --keystore=<path-to-keystore>
```

### Create a multisig key

The following example creates a multisig key with 3 sub-keys, and specify the minimum number of signatures as 2. The tx could be broadcast only when the number of signatures is greater than or equal to 2.

```bash
iriscli keys add <multisig-keyname> --multisig-threshold=2 --multisig=<signer-keyname-1>,<signer-keyname-2>,<signer-keyname-3>
```

:::tip
`<signer-keyname>` can be the type of "local/offline/ledger", but not "multi" type.

If you don't have all the permission of sub-keys, you can ask for the pubkeys to create the offline keys first, then you will be able to create the multisig key.

Offline key can be created by "iriscli keys add --pubkey".
:::

How to use multisig key to sign and broadcast a transaction,  please refer to [multisign](tx.md#iriscli-tx-multisign)

## iriscli keys list

List all the keys stored by this key manager along with their associated name, type, address and pubkey.

### List all keys

```bash
iriscli keys list
```

## iriscli keys show

Get details of a local key.

```bash
iriscli keys show <key-name> <flags>
```

**Flags:**

| Name, shorthand      | Default | Description                                         | Required |
| -------------------- | ------- | --------------------------------------------------- | -------- |
| --address            |         | Output the address only (overrides --output)        |          |
| --bech               | acc     | The Bech32 prefix encoding for a key (acc/val/cons) |          |
| --help, -h           |         | help for show                                       |          |
| --multisig-threshold | 1       | K out of N required signatures                      |          |
| --pubkey             |         | Output the public key only (overrides --output)     |          |

### Get details of a local key

```bash
iriscli keys show MyKey
```

The following infos will be shown:

```bash
NAME:    TYPE:    ADDRESS:                                      PUBKEY:
MyKey    local    iaa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zak83e7nf    iap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskzc7exa
```

### Get validator operator address

If an address has been bonded to be a validator operator (which the address you used to create a validator), then you can use `--bech val` to get the operator's address prefixed by `iva` and the pubkey prefixed by `ivp`:

```bash
iriscli keys show MyKey --bech val
```

Example Output:

```bash
NAME:    TYPE:    ADDRESS:                                      PUBKEY:
MyKey    local    iva12nda6xwpmp000jghyneazh4kkgl2tnzyx7trze    ivp1addwnpepqfw52vyzt9xgshxmw7vgpfqrey30668g36f9z837kj9dy68kn2wxqm8gtmk
```

## iriscli keys export

Export the keystore of a key to a json file

```bash
iriscli keys export <key-name> <flags>
```

**Flags:**

| Name, shorthand | Default | Description          | Required |
| --------------- | ------- | -------------------- | -------- |
| --output-file   |         | The path of keystore |          |

### Export keystore

```bash
iriscli keys export Mykey --output-file=<path-to-keystore>
```

## iriscli keys delete

Delete a local key by the given name.

```bash
iriscli keys delete <key-name> <flags>
```

**Flags:**

| Name, shorthand | Default | Description                                                             | Required |
| --------------- | ------- | ----------------------------------------------------------------------- | -------- |
| --force, -f     | false   | Remove the key unconditionally without asking for the passphrase        |          |
| --yes, -y       | false   | Skip confirmation prompt when deleting offline or ledger key references |          |

### Delete a local key

```bash
iriscli keys delete MyKey
```

## iriscli keys update

Change the password of a key, used to protect the private key.

### Change the password of a local key

```bash
iriscli keys update MyKey
```

## iriscli keys mnemonic

Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. To pass your own entropy, use `unsafe-entropy` mode.

```bash
iriscli keys mnemonic <flags>
```

**Flags:**

| Name, shorthand  | Default | Description                                                                   | Required |
| ---------------- | ------- | ----------------------------------------------------------------------------- | -------- |
| --unsafe-entropy |         | Prompt the user to supply their own entropy, instead of relying on the system |          |

### Create a bip39 mnemonic

```bash
iriscli keys mnemonic
```

You'll get a bip39 mnemonic with 24 words, e.g.:

```bash
police possible oval milk network indicate usual blossom spring wasp taste canal announce purpose rib mind river pet brown web response sting remain airport
```

## iriscli keys new

:::warning
**Deprecated**
:::

Derive a new private key using an interactive command that will prompt you for each input.

Optionally specify a bip39 mnemonic, a bip39 passphrase to further secure the mnemonic, and a bip32 HD path to derive a specific account. The key will be stored under the given name and encrypted with the given password. The only input that is required is the encryption password.

```bash
iriscli keys new <key-name> <flags>
```

**Flags:**

| Name, shorthand | Default         | Description                                                     | Required |
| --------------- | --------------- | --------------------------------------------------------------- | -------- |
| --bip44-path    | 44'/118'/0'/0/0 | BIP44 path from which to derive a private key                   |          |
| --default       |                 | Skip the prompts and just use the default values for everything |          |
| --help, -h      |                 | Help for new                                                    |          |
| --ledger        |                 | Store a local reference to a private key on a Ledger device     |          |

### Create a new key by the specified method

```bash
iriscli keys new MyKey
```

You'll be asked to enter your bip44 path, default is 44'/118'/0'/0/0.

```bash
> -------------------------------------
> Enter your bip44 path. Default is 44'/118'/0'/0/0
```

Then you'll be asked to enter your bip39 mnemonic, or hit enter to generate one.

```bash
> Enter your bip39 mnemonic, or hit enter to generate one.
```

You can hit enter to generate bip39 mnemonic, then a new hint will be show to ask you to enter bip39 passphrase.

```bash
> -------------------------------------
> Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed
> Most users should just hit enter to use the default, ""
```

Also you can hit enter to skip it, then you'll receive a hint to enter a password.

```bash
> -------------------------------------
> Enter a passphrase to encrypt your key to disk:
> Repeat the passphrase:
```

After that, you're done with creating a new key.
