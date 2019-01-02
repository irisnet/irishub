#!/bin/bash

# The default values in this script are for mainnet
# Edit the script to set environment variables for your own network

# Network type: testnet, mainnet
export NetworkType="mainnet"

# Bech32 prefixes
export Bech32PrefixAccAddr="iaa"
export Bech32PrefixAccPub="iap"
export Bech32PrefixValAddr="iva"
export Bech32PrefixValPub="ivp"
export Bech32PrefixConsAddr="ica"
export Bech32PrefixConsPub="icp"

# Disable invariant checking: "no"
# Panic on invariant failure: "panic"
# Print error message on invariant failure: "error"
export InvariantLevel="error"