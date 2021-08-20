---
order: 4
---

# Join The Testnet

After IRIS Hub 1.0 upgrade of mainnet, **Nyancat Testnet** starts to operate as a stable application testnet which has the same version as the mainnet, so that the service providers of IRISnet can develop their apps on or interact with IRIShub without running a node or lcd instance.

## Public Endpoints

- GRPC: 35.234.10.84:9090
- RPC: http://35.234.10.84:26657/
- REST: http://35.234.10.84:1317/swagger/



## Run a Full Node

### Start node from genesis
::tip 
You must use Irishub [v1.1.1](https://github.com/irisnet/irishub/releases/tag/v1.1.1)[ ](https://github.com/irisnet/irishub/releases/tag/v1.0.1) to initialize your node::

```bash
# init node
iris init <moniker> --chain-id=nyancat-8

# download public config.toml and genesis.json
curl -o ~/.iris/config/config.toml https://github.com/irisnet/testnets/blob/master/nyancat/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/testnets/master/nyancat/config/genesis.json

# Start the node (also running in the background, such as nohup or systemd)
iris start
```



## Faucet

Welcome to get test tokens in our [testnet faucet channel](https://discord.gg/Z6PXeTb5Mt) 

Usage: in [nyancat-faucet channel](https://discord.gg/Z6PXeTb5Mt), type "$faucet " + your address on Nyancat testnet, to get test tokens (NYAN) every 24 hours.

## Explorer

<https://nyancat.iobscan.io/>

## Community

Welcome to discuss in our [nyancat testnet channel](https://discord.gg/9cSt7MX2fn) 

## Faucet

Welcome to ask for the test tokens in our [testnet channels](https://discord.gg/9cSt7MX2fn)

## Explorer

<https://nyancat.iobscan.io/>
