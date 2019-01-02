#!/bin/bash

# Edit the script to set compile environment variables for your own application

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

export NetworkType="mainnet"
