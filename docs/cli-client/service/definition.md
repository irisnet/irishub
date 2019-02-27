# iriscli service definition

## Description

Query service definition

## Usage

```
iriscli service definition [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --def-chain-id  |                            | [string] the ID of the blockchain defined of the service            | Yes      |
| --service-name  |                            | [string] service name                                               | Yes      |
| --help, -h      |                            | help for definition                                                 |          |

## Examples

### Query a service definition

```shell
iriscli service definition --def-chain-id=<chain-id> --service-name=test-service
```

After that, you will get detail info for the service definition which has the specfied define chain id and service name.

```json
{
  "SvcDef": {
    "name": "test-service",
    "chain_id": "test",
    "description": "service-description",
    "tags": [
      "tag1",
      "tag2"
    ],
    "author": "iaa1ydhmma8l4m9dygsh7l08fgrwka6yczs0se0tvs",
    "author_description": "author-description",
    "idl_content": "syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n"
  },
  "methods": [
    {
      "id": "1",
      "name": "SayHello",
      "description": "sayHello",
      "output_privacy": "NoPrivacy",
      "output_cached": "NoCached"
    }
  ]
}
```