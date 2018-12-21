# iriscli service define 

## Description

Create a new service definition

## Usage

```
iriscli service define [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --service-description |                         | [string] service description                                                                                                                          |          |
| --author-description  |                         | [string] service author description                                                                                                                   |          |
| --service-name        |                         | [string] service name                                                                                                                                 |   Yes    |
| --tags                |                         | [strings] service tags                                                                                                                                |          |
| --idl-content         |                         | [string] content of service interface description language                                                                                            |          |
| --file                |                         | [string] path of file which contains service interface description language                                                                           |          |
| -h, --help            |                         | help for define                                                                                                                                       |          |

## Examples

### define a service
```shell
iriscli service define --chain-id=test  --from=node0 --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```
Idl-content can be replaced by file if the file item is not empty.  [Example of IDL content](#idl-content-example).

After that, you're done with defining a new service.

```txt
Committed at block 539 (tx hash: 9ED8B36F8DDA7745BF03E0F5271E55B6D0BC34B373BFCDB6B5BC78C502DAE032, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7604,
   "codespace": "",
   "tags": {
     "action": "service_define"
   }
 })
```

### IDL content example
* idl-content example

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)