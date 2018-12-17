# iriscli service respond 

## Description

Respond a service method invocation

## Usage

```
iriscli service respond [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] the ID of the blockchain that the service invocation initiated                                                                                              |  Yes     |
| --request-id          |                         | [string] the ID of the service invocation                                                                                                                                |  Yes     |
| --response-data       |                         | [string] hex encoded response data of a service invocation                                                                       |         |
| -h, --help            |                         | help for respond                                                                                                                                         |          |

## Examples

### Respond to a service invocation 
```shell
iriscli service respond --chain-id=test-irishub --from=node0 --fee=0.004iris --request-chain-id=test-irishub --request-id=230-130-0 --response-data=abcd
```

After that, you're done with responding to a service invocation.

```json
{
   "tags": {
     "action": "service-call",
     "completeConsumedTxFee-iris-atto": "\"162880000000000\"",
     "consumer": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
     "provider": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
     "request-id": "230-130-0"
   }
 }
```

