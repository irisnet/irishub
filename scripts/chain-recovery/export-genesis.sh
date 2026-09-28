#!/usr/bin/env bash

set -euo pipefail

IRIS_BIN=/path/to/verified/iris
OLD_HOME=/path/to/stopped-node-copy
OUT_DIR=/path/to/export-output

NEW_CHAIN_ID=irishub-2
for output in "$OUT_DIR/default" "$OUT_DIR/partial.json" "$OUT_DIR/genesis.json"; do
  if [[ -e "$output" ]]; then
    printf 'refusing to overwrite existing path: %s\n' "$output" >&2
    exit 1
  fi
done
mkdir -p "$OUT_DIR"

"$IRIS_BIN" init recovery \
  --chain-id "$NEW_CHAIN_ID" \
  --home "$OUT_DIR/default" >/dev/null

MODULES=$(jq -r '
  .app_state | keys -
  ["07-tendermint","ibc","transfer","interchainaccounts","nonfungibletokentransfer"]
  | join(",")
' "$OUT_DIR/default/config/genesis.json")

"$IRIS_BIN" export \
  --home "$OLD_HOME" \
  --height 37242247 \
  --for-zero-height=false \
  --modules-to-export "$MODULES" \
  --output-document "$OUT_DIR/partial.json"

jq --slurpfile defaults "$OUT_DIR/default/config/genesis.json" \
    --arg chain_id "$NEW_CHAIN_ID" '
  .chain_id = $chain_id
  | reduce ["07-tendermint","ibc","transfer","interchainaccounts",
            "nonfungibletokentransfer"][] as $m
      (. ; .app_state[$m] = $defaults[0].app_state[$m])
' "$OUT_DIR/partial.json" > "$OUT_DIR/genesis.json"

"$IRIS_BIN" genesis validate "$OUT_DIR/genesis.json"

jq -e --arg chain_id "$NEW_CHAIN_ID" \
  '.initial_height == 37242248 and .chain_id == $chain_id' \
  "$OUT_DIR/genesis.json" >/dev/null

EXPORTED_SUPPLY=$(jq -r '
  .app_state.bank.supply[]
  | select(.denom == "uiris")
  | .amount
' "$OUT_DIR/genesis.json")
test "$EXPORTED_SUPPLY" = "2156713998266827"

printf 'genesis: %s\nuiris supply: %s\n' \
  "$OUT_DIR/genesis.json" "$EXPORTED_SUPPLY"
