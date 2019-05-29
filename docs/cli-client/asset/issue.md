# iriscli asset issue

## Description

This command is used to issue asset in IRIS Hub.

## Usage

issue 100000000000 iris
```
iriscli asset issue [flags]
```


## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                                                                                                                                                                                                                                           |
| ---------------- | --------- | ------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------                                                                                                                                     |
| --family         | byte      | false       |  00                | 00 - fungible;01 - non-fungible                                                                                                                                                                                                                                                  |
| --name           | string    | true        |                    | Name of the newly issued asset, limited to 32 unicode(english character , number, _) characters,  e.g. "IRISnet"                                                                                                                                                                 |
| --symbol         | string    | true        |                    | The length of the string for representing this asset is between 3 and 6 alphanumeric characters and is case insensitive                                                                                                                                                          |
| --source         | string    | false       |  00                | Reserved - 00 (native); 01 (external); Gateway IDs                                                                                                                                                                                                                               |
| --initial-supply | uint64    | true        |                    | The initial supply for this asset. The amount before boosting should not exceed 100 billion. The amount should be positive integer                                                                                                                                               |
| --max-supply     | uint64    | true        |  1000000000000     | The hard cap of this asset, total supply can not exceed max supply. The amount should be positive integer                                                                                                                                                                        |
| --decimal        | uint8     | false       |  0                 | The asset can have a maximum of 18 digits of decimal                                                                                                                                                                                                                             |
| --mintable       | boolean   | false       |  false             | Whether this asset could be minted(increased) after the initial issuing                                                                                                                                                                                                          |
| --operators      | []Address | false       |  []                | Operators have all permissions of this asset except transfering owner, but if the owner lost his private key and there are more than one operators, then the operators can transfer the asset owner by sending an asset transfer tx signed with no less than 2/3 of the operators|




## Examples

### Issue asset

```
iriscli asset issue iris --family=00 --name=IRISnet --symbol=iris --source=00 --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintalbe=false  --operators=<account address A>,<account address B>

```

#### TODO:After that, you will get the detail info for the account.
