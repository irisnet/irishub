# Asset User Guide

## Introduction

This specification describes asset managerment on IRISHub. Anyone could issue a new asset on IRISHub, or propose pegging an existing asset from any other blockchains via On-Chain Governance.

## Concepts

### Gateways

A gateway is a trusted party that facilitates moving value into and out of the IRIS Network. Gateways are basically equivalent to the standard exchange model where you depend on the solvency of the exchange to be able to redeem your coins. Generally gateways issue assets prefixed with their symbol, like GDEX, OPEN, and so on. These assets are backed 100% by the real BTC or ETH or any other coin that people deposit with the gateways.

### Assets

#### User Issued Assets (UIAs)

IRISHub allows individuals and companies to create and issue their own tokens for anything they can imagine. The potential use cases for so called user-issued assets (UIA) are innumerable. On the one hand, UIAs can be used as simple event tickets deposited on the customers mobile phone to pass the entrance of a concert. On the other hand, they can be used for crowd funding, ownership tracking or even to sell equity of a company in form of stock.
All you need to do in order to create a new UIA is to execute an one-line command, define your preferred parameters for your coin, such as supply, symbol, description etc. From that point on, you can issue some of your coins to whomever you want just like sending iris.
As the owner of that asset, you don’t need to take care of any the technical details of blockchain technology, such as distributed consensus algorithms, blockchain development or integration. You don’t even need to run any mining equipment or servers, at all.

#### Market Pegged Assets (MPAs)

Instead of creating an UIA where the full control over supply is in the hands of the issuer, we can also create a Market Pegged Asset (MPA) and let the market deal with demand and supply.
The only way to create a MPA is to create a proposal on IRISHub Governance [*TODO*], however the top 20 of CMC will be pre-configured in system.

## Actions

[*TODO*]

## Reference

[BitShares Documentation](https://how.bitshares.works/en/master/index.html)

[Benefits of Becoming a BitShares Gateway](http://bytemaster.github.io/update/2014/12/18/Benefits-of-Being-a-BitShares-Gateway/)
