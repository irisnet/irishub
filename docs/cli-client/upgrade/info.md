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

Then it will show 

```
{
"current_proposal_id": "0",
"current_proposal_accept_height": "-1",
"version": {
"Id": "0",
"ProposalID": "0",
"Start": "0",
"ModuleList": [
{
"Start": "0",
"End": "9223372036854775807",
"Handler": "bank",
"Store": [
"acc"
]
},
{
"Start": "0",
"End": "9223372036854775807",
"Handler": "stake",
"Store": [
"stake",
"acc",
"mint",
"distr"
]
},
.......
]
}
}
```
