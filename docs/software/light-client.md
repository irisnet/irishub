# IRISLCD

## Introduction

An IRISLCD node is a light node of IRISHUB. Unlike IRISHUB full node, it won't store all blocks and execute all transactions, which means it only requires minimal bandwidth, computing and storage resource. In distrust mode, it will track the evolution of validator set change and require full nodes to return consensus proof and merkle proof. Unless validators with more than 2/3 voting power do byzantine behavior, then IRISLCD proof verification algorithm can detect all potential malicious data, which means an IRISLCD node can provide the same security as full nodes.

The default home folder of irislcd is `$HOME/.irislcd`. Once an IRISLCD is started, it will create two directories: `keys` and `trust-base.db`.The keys store db locates in `keys`. `trust-base.db` stores all trusted validator set and other verification related files.

When IRISLCD is started in distrust mode, it will check whether `trust-base.db` is empty. If true, it will fetch the latest block as its trust basis and save it under `trust-base.db`. The IRISLCD node always trust the basis. All query proof will be verified based on the trust basis, which means IRISLCD can only verify the proof on later heights. If you want to query transactions or blocks on lower heights, please start IRISLCD in trust mode. For detailed proof verification algorithm please refer to [tendermint lite](https://github.com/tendermint/tendermint/blob/master/docs/tendermint-core/light-client-protocol.md).

## Usage

For how to start IRISLCD and how to access these REST APIs, please refer to [how to use light-client](../light-client/README.md).
