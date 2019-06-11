# Asset User Guide

## Introduction

This specification describes asset managerment on IRISHub. Anyone could issue a new asset on IRISHub, or propose pegging an existing asset from any other blockchains via On-Chain Governance.

## Concepts

### Assets

#### Native Assets

IRISHub allows individuals and companies to create and issue their own tokens for anything they can imagine. The potential use cases are innumerable. On the one hand, native assets can be used as simple event tickets deposited on the customers mobile phone to pass the entrance of a concert. On the other hand, they can be used for crowd funding, ownership tracking or even to sell equity of a company in form of stock.
All you need to do in order to create a new `native asset` is to execute an one-line command, define your preferred parameters for your coin, such as supply, symbol, description etc. From that point on, you can issue some of your coins to whomever you want just like sending iris.
As the owner of that asset, you don’t need to take care of any the technical details of blockchain technology, such as distributed consensus algorithms, blockchain development or integration. You don’t even need to run any mining equipment or servers, at all.

#### External Assets

Instead of creating a `native asset` where the full control over supply is in the hands of the issuer, we can also create an `external asset` which is already exists on any other blockchains and let the market deal with demand and supply.
The only way to create an `external asset` is to create an `add-asset`[*TODO*] proposal on IRISHub Governance, however the top 20 of CMC will be pre-configured in system.

### Gateways

A gateway is a trusted party that facilitates moving value into and out of the IRIS Network. Gateways are basically equivalent to the standard exchange model where you depend on the solvency of the exchange to be able to redeem your coins. Generally gateways issue [native assets](#Native-Assets) prefixed with their symbol, like GDEX, OPEN, and so on. These assets are backed 100% by the real BTC or ETH or any other coin that people deposit with the gateways.

## Actions

[*TODO*]
