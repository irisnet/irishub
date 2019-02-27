# iriscli keys list

## Description

Return a list of all public keys stored by this key manager
along with their associated name and address.

## Usage

```
iriscli keys list [flags]
```

## Flags

| Name, shorthand | Default   | Description                                                  | Required |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | help for list                                                |          |

## Examples

### List all public keys 

```shell
iriscli keys list
```

You'll get all the local public keys with 'address' and 'pubkey' element.

```txt
NAME:	TYPE:	ADDRESS:						            PUBKEY:
abc  	local	iaa1va2eu9qhwn5fx58kvl87x05ee4qrgh44yd8teh	iap1addwnpepq02r0hts0yjhp4rsal627s2lqk4agy2g6tek5g9yq2tfrmkkehee2td75cs
test	local	iaa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zakl7lxn5	iap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskww57lw
```
