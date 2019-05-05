# iriscli keys mnemonic

## Description

Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. To pass your own entropy, use `unsafe-entropy` mode.

## Usage

```
iriscli keys mnemonic <flags>
```

## Flags

| Name, shorthand  | Default   | Description                                                                   | Required |
| ---------------- | --------- | ----------------------------------------------------------------------------- | -------- |
| --help, -h       |           | help for mnemonic                                                             |          |
| --unsafe-entropy |           | Prompt the user to supply their own entropy, instead of relying on the system |          |

## Examples

### Create a bip39 mnemonic

```shell
iriscli keys mnemonic
```

You'll get a bip39 mnemonic with 24 words.

```txt
police possible oval milk network indicate usual blossom spring wasp taste canal announce purpose rib mind river pet brown web response sting remain airport
```

### Use `unsafe-entropy` mode.

```shell
root@ubuntu16:~# iriscli keys mnemonic --unsafe-entropy

WARNING: Generate at least 256-bits of entropy and enter the results here:
<input_your_own_entropy_string>
Input length: 128 [y/n]:y
-------------------------------------
wine hire tongue weasel air puzzle claim pole curtain taste box learn exchange where become inside blur tragic suffer fruit hole transfer race unit
```