# FAQ

1. What is IRISnet?

IRIS network is named after the Greek goddess Iris, said to be the personification of the rainbow and the faithful messenger between heaven and humanity. IRISnet is designed to be the foundation for next generation distributed business applications. IRISnet will enhance Interblockchain Communication(IBC)protocol of Cosmos to introduce a service-oriented infrastructure into ecosystem to support more efficient building distributed business application chains. It will make services interoperable across blockchains: public, consortium & legacy systems.

2. IRISnet and Cosmos

Cosmos IBC defines a protocol for transferring values from an account on one chain to an account on another chain. The IRIS network designs new semantics to allow cross-chain computation to be invoked by leveraging IBC.The IRIS network could provide the service infrastructure for handing and coordinating on-chain transaction processing with off-chain data processing and business logic execution


3. What is IRIS token?

IRIS token is the native toke in IRISnet. It has two main usages:

* Staking token: IRIS token holders could stake or delegate some IRIS to become a validator candidate
* Fee token: IRIS could be used to pay for network fee and service fee.

iris precision: 10^18



3. Server Configuration

Here is the recommanded configuration of the server.
* No. of CPUs: 2
* Memory: 4GB
* Disk: 40GB SSD
* OS: Ubuntu 18.04 LTS/16.04 LTS
* Allow all incoming connections from TCP port 26656 and 26657


4. What is a Validator?


The IRISHub is based on a consensus engine called Tendermint. It relies on a set of validators to secure the network. The role of validators is to run a full-node and participate in consensus by broadcasting votes which contain cryptographic signatures signed by their private key. Validators commit new blocks in the blockchain and receive revenue in exchange for their work. They must also participate in governance by voting on proposals. Validators are weighted according to their total stake.  

The reward for a validator is 3-fold:

* Block provision
* Transaction fee
* Commission


5. What is a Delegator?

Not every iris token holder is eligible to become a validator. Any iris token holder could choose to delegate their own iris token to one or more validators. In this way, they could still earn block provision and transaction fees. At the same time, some percent of commission needs to be paid to their validator.


6. How to understand different address formats in IRISnet?

Please read this [doc](https://github.com/irisnet/testnets/blob/master/fuxi/docs/Bech32%20on%20IRISnet.md) to understand the address format in IRISnet. 

7. How to understand the notion of gas&fee in IRISHub?

This is the tech spec for fee&gas is [here](https://github.com/irisnet/irishub/blob/d1d20826da2112a53c6a0ce45e0263237c549089/docs/modules/fee-token/feeToken.md)


