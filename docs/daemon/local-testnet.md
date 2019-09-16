# Local Testnet

For testing or developing purpose, you may want to setup a local testnet.

## Single Node Testnet

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
iris init --chain-id=testing --moniker=testing

# Create a key to hold your validator account
iriscli keys add MyValidator

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some iris
# which is the only staking coin on IRISnet
iris add-genesis-account $(iriscli keys show MyValidator --address) 100000000iris

# Generate the transaction that creates your validator
# The gentxs are stored in ~/.iris/config/gentx/
iris gentx --name MyValidator

# Add the generated bonding transactions to the genesis file
iris collect-gentxs

# Now its ready to start `iris`
iris start
```

## Multiple Nodes Testnet

TODO
