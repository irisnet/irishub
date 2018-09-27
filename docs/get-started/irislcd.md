# What is Irislcd

An irislcd node is a REST server which can connect to any full nodes and provide a set of rest APIs. By these APIs, users can send transactions and query blockchain data. Irislcd can verify the proof of query result. So it can provide the same security as a full node with the minimal requirements on bandwidth, computing and storage resource. Besides, it also provides swagger-ui which presents detailed description about what APIs it provides and how to use the them. 

## Irislcd options

To start a irislcd, we need to specify the following parameters:

| Parameter       | Type      | Default                 | Required | Description                                          |
| --------------- | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| chain-id        | string    | null                    | true     | Chain ID of Tendermint node |
| home            | string    | "$HOME/.irislcd"        | false    | Directory for config and data, such as key and checkpoint |
| node            | string    | "tcp://localhost:26657" | false    | Full node to connect to |
| laddr           | string    | "tcp://localhost:1317"  | false    | Address for server to listen on |
| trust-node      | bool      | false                   | false    | Trust connected  full nodes (Don't verify proofs for responses) |
| max-open        | int       | 1000                    | false    | The number of maximum open connections |
| cors            | string    | ""                      | false    |Set the domains that can make CORS requests |

## Start Irislcd

Sample command to start irislcd:
```
irislcd start --chain-id=<chain-id>
```
Please visit the following url with in your internet explorer to open Swagger-UI:
```
http://localhost:1317/swagger-ui/
```
Execute the following command to print the irislcd version.
```
irislcd version
```

When the connected full node is trusted, then the proof is not necessary, so you can run irislcd with trust-node option:
```
irislcd start --chain-id=<chain-id> --trust-node