#!/bin/bash

# Load shell variables
. ./network/hermes/variables.sh

sleep 2

# Start the hermes relayer in multi-paths mode
echo "Starting hermes relayer..."
$HERMES_BINARY --config $CONFIG_DIR start
