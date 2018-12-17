# iriscli upgrade info

## Description

Query the information of software version and upgrade module.

## Usage

```
iriscli upgrade info
```

Print help messages:

```
iriscli upgrade info --help
```
## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |

## Example

Query the current version information. 

```
iriscli upgrade info 
```

Then it will show the current protocol information and the protocol information is currently being prepared for upgrade.

```
{
"version": {
"ProposalID": "1",
"Success": true,
"Protocol": {
"version": "1",
"software": "https://github.com/irisnet/irishub/tree/v0.7.0",
"height": "30"
}

},
"upgrade_config": {
"ProposalID": "3",
"Definition": {
"version": "2",
"software": "https://github.com/irisnet/irishub/tree/v0.9.0",
"height": "80"
}
}
}
```
