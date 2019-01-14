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
"version": "0",
"software": "https://github.com/irisnet/irishub/tree/v0.10.0",
"height": "1"
}

},
"upgrade_config": {
"ProposalID": "3",
"Definition": {
"version": "1",
"software": "https://github.com/irisnet/irishub/tree/v0.10.1",
"height": "8000"
}
}
}
```
