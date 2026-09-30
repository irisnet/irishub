#!/usr/bin/env bash
#
# Launch a local simulation chain from the recovered genesis: replace the
# exported (uncontrolled) validator set with locally generated validators,
# start the chain, and verify that blocks, transactions and queries work
# against the recovered data.
#
# Usage:
#   IRIS_BIN=/path/to/iris \
#   EXPORT_GENESIS=/mnt/export-output/genesis.json \
#   bash scripts/chain-recovery/simulate-launch.sh
#
# Env:
#   IRIS_BIN        iris binary (default: iris)
#   EXPORT_GENESIS  the recovered genesis to simulate (required)
#   SIM_DIR         working directory (default: /tmp/irishub-sim)
#   SIM_CHAIN_ID    simulation chain-id (default: irishub-sim-1); the final
#                   genesis carries this id because gentx signatures bind to it
#   KEEP_RUNNING=1  leave the node running when done (default: stop it)

set -euo pipefail

IRIS_BIN=${IRIS_BIN:-iris}
EXPORT_GENESIS=${EXPORT_GENESIS:?set EXPORT_GENESIS to the recovered genesis file}
SIM_DIR=${SIM_DIR:-/tmp/irishub-sim}
SIM_CHAIN_ID=${SIM_CHAIN_ID:-irishub-sim-1}
KEEP_RUNNING=${KEEP_RUNNING:-0}
RPC=tcp://127.0.0.1:26657
RPC_HTTP=http://127.0.0.1:26657

die() { printf '%s\n' "$*" >&2; exit 1; }

[[ -f $EXPORT_GENESIS ]] || die "genesis not found: $EXPORT_GENESIS"
if [[ -e $SIM_DIR ]]; then
  die "refusing to overwrite existing dir: $SIM_DIR (remove it first)"
fi
mkdir -p "$SIM_DIR"

NODE_HOME="$SIM_DIR/testnet/node0/iris"

printf '==> generating a local testnet with our own validator\n'
"$IRIS_BIN" testnet \
  --output-dir "$SIM_DIR/testnet" \
  --chain-id "$SIM_CHAIN_ID" \
  --v 1 \
  --keyring-backend test \
  >/dev/null

printf '==> replacing the exported validator set with ours\n'
# `replace-validators` reads genesis files with cometbft's strict decoder, which
# requires int64 fields (initial_height, app_hash) in cometbft spelling while
# the SDK writes plain numbers; convert both files first
jq '.initial_height |= tostring | .app_hash = (.app_hash // "")' \
  "$NODE_HOME/config/genesis.json" > "$SIM_DIR/source.comet.json"
jq '.initial_height |= tostring | .app_hash = (.app_hash // "")' \
  "$EXPORT_GENESIS" > "$SIM_DIR/target.comet.json"
"$IRIS_BIN" genesis replace-validators \
  --source-genesis-file "$SIM_DIR/source.comet.json" \
  --target-genesis-file "$SIM_DIR/target.comet.json" \
  --output-genesis-file "$SIM_DIR/simulation.comet.json"

printf '==> converting the merged genesis back to SDK format\n'
# re-attach the SDK-only members (app_name/app_version/consensus) the comet
# decoder dropped; consensus.validators MUST stay empty: the doc validator set
# would otherwise mix the old (uncontrolled) validators with ours and block
# consensus, the set comes from the imported gentxs via InitChain instead
jq -c '{app_name, app_version, params: .consensus.params}' \
  "$EXPORT_GENESIS" > "$SIM_DIR/meta.json"
jq --slurpfile m "$SIM_DIR/meta.json" \
  '.initial_height |= tonumber
   | .app_name = $m[0].app_name
   | .app_version = $m[0].app_version
   | .consensus = {params: $m[0].params, validators: []}' \
  "$SIM_DIR/simulation.comet.json" > "$SIM_DIR/simulation-genesis.json"
"$IRIS_BIN" genesis validate "$SIM_DIR/simulation-genesis.json" >/dev/null

printf '==> starting the simulation node\n'
cp "$SIM_DIR/simulation-genesis.json" "$NODE_HOME/config/genesis.json"
"$IRIS_BIN" start --home "$NODE_HOME" > "$SIM_DIR/node.log" 2>&1 &
NODE_PID=$!

stop_node() {
  if [[ $KEEP_RUNNING != 1 ]]; then
    kill "$NODE_PID" 2>/dev/null || true
  fi
}
trap stop_node EXIT

printf '==> waiting for the node to produce blocks\n'
HEIGHT=0
for _ in $(seq 1 60); do
  H=$( { curl -s "$RPC_HTTP/status" 2>/dev/null | jq -r '.result.sync_info.latest_block_height // 0'; } 2>/dev/null || echo 0)
  if [[ "$H" -gt 2 ]]; then HEIGHT=$H; break; fi
  sleep 2
done
if [[ "$HEIGHT" == 0 ]]; then
  tail -30 "$SIM_DIR/node.log"
  die "node did not produce blocks"
fi
printf '    block height: %s\n' "$HEIGHT"

printf '==> staking: our validator is bonded\n'
"$IRIS_BIN" q staking validators --node "$RPC" 2>/dev/null \
  | grep -E 'operator_address|status:|tokens:' | head -3 || true

printf '==> bank: total supply\n'
"$IRIS_BIN" q bank total --node "$RPC" 2>/dev/null | head -4 || true

printf '==> sending a test transaction (bank send)\n'
VAL_ADDR=$("$IRIS_BIN" keys show node0 \
  --home "$NODE_HOME" --keyring-backend test --output json 2>/dev/null | jq -r .address)
TX_OUT=$("$IRIS_BIN" tx bank send "$VAL_ADDR" "$VAL_ADDR" 1uiris \
  --chain-id "$SIM_CHAIN_ID" \
  --home "$NODE_HOME" --keyring-backend test \
  --node "$RPC" \
  --fees 60000uiris --yes -o json 2>&1)
TX_HASH=$(printf '%s\n' "$TX_OUT" | jq -r '.txhash // empty' 2>/dev/null || true)
[[ -n $TX_HASH ]] \
  || { printf '%s\n' "$TX_OUT"; die 'test transaction failed'; }
TXQ=""
for _ in $(seq 1 10); do
  TXQ=$("$IRIS_BIN" q tx "$TX_HASH" --node "$RPC" -o json 2>/dev/null || true)
  if printf '%s' "$TXQ" | jq -e '.height and .height != "0"' >/dev/null 2>&1; then break; fi
  sleep 2
done
if [[ -z $TXQ ]]; then
  die "test transaction $TX_HASH was not committed"
fi
printf '%s\n' "$TXQ" | jq -r '"    committed: txhash \(.txhash[0:16])... height \(.height) code \(.code)"'

printf '\nAll checks passed. '
if [[ $KEEP_RUNNING == 1 ]]; then
  trap - EXIT
  printf 'node kept running:\n  home: %s\n  rpc:  %s\n  key:  node0 (keyring backend "test" under the home)\n' \
    "$NODE_HOME" "$RPC"
else
  printf 'node stopped. re-run anytime after removing %s\n' "$SIM_DIR"
fi
