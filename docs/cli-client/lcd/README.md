# What is Irislcd

An IRISLCD node is a REST server which can connect to any full nodes and provide a set of rest APIs. By these APIs, users can send transactions and query blockchain data. Irislcd can verify the proof of query result. So it can provide the same security as a full node with the minimal requirements on bandwidth, computing and storage resource. Besides, it also provides swagger-ui which presents detailed description about what APIs it provides and how to use them. 

## Irislcd usage

Irislcd has two subcommands:

| subcommand      | Description                 | Example command |
| --------------- | --------------------------- | --------------- |
| version         | Print the IRISLCD version   | IRISLCD version |
| start           | Start a IRISLCD node        | IRISLCD start --chain-id=<chain-id> |

`start` subcommand has these options:

| Parameter       | Type      | Default                 | Required | Description                                          |
| --------------- | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| chain-id        | string    | null                    | true     | Chain ID of Tendermint node |
| home            | string    | "$HOME/.irislcd"        | false    | Directory for config and data, such as key and checkpoint |
| node            | string    | "tcp://localhost:26657" | false    | Full node to connect to |
| laddr           | string    | "tcp://localhost:1317"  | false    | Address for server to listen on |
| trust-node      | bool      | false                   | false    | Trust connected  full nodes (Don't verify proofs for responses) |
| max-open        | int       | 1000                    | false    | The number of maximum open connections |
| cors            | string    | ""                      | false    | Set the domains that can make CORS requests |

## Sample commands

1. When the connected full node is trusted, then the proof is not necessary, so you can run IRISLCD with trust-node option:
```bash
irislcd start --chain-id=<chain-id> --trust-node
```

2. If you want to access your IRISLCD in remote machine, you have to specify `--laddr`, for instance:
```bash
irislcd start --chain-id=<chain-id> --laddr=tcp://0.0.0.0:1317
```