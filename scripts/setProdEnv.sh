#!/bin/bash

# The default values in this script are for mainnet
# Edit the script to set environment variables for your own network

# Network type: testnet, mainnet
export NetworkType="mainnet"

# Disable invariant checking: "no"
# Panic on invariant failure: "panic"
# Print error message on invariant failure: "error"
export InvariantLevel="error"
