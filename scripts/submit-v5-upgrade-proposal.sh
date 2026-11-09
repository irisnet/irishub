#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Generate and optionally submit the IRIS Hub v5.0.0 software-upgrade proposal.

Usage:
  KEY_NAME=<key> UPGRADE_HEIGHT=<height> GAS_PRICES=<price> \
    scripts/submit-v5-upgrade-proposal.sh [--submit]

The default mode queries the chain, writes the proposal artifacts, and generates
an unsigned transaction. It does not broadcast. Pass --submit to sign and
broadcast the proposal after the unsigned transaction has been generated.

Required environment variables:
  KEY_NAME          Key name or signer address used by iris tx
  UPGRADE_HEIGHT    Future block height at which v5.0.0 will execute
  GAS_PRICES        Gas price accepted by the target node, for example 0.3uiris

Optional environment variables:
  IRIS_BIN          iris executable (default: iris)
  CHAIN_ID          Target chain ID (default: irishub-2)
  NODE              CometBFT RPC address (default: tcp://127.0.0.1:26657)
  IRIS_HOME         iris client home (default: $HOME/.iris)
  KEYRING_BACKEND   Keyring backend (default: os)
  UPGRADE_INFO      Upgrade Plan info field (default: empty)
  OUTPUT_DIR        Artifact directory (default: ./irishub-v5-upgrade-proposal)

UPGRADE_INFO is intentionally empty by default. Set it only after official
v5.0.0 binaries, checksums, and download instructions have been published.
This script is intended for a normal single-signer account; multisig proposals
must use the corresponding offline multisig signing flow.
EOF
}

submit=false
case "${1:-}" in
  "") ;;
  --submit) submit=true ;;
  -h|--help)
    usage
    exit 0
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
if [[ $# -gt 1 ]]; then
  usage >&2
  exit 2
fi

IRIS_BIN="${IRIS_BIN:-iris}"
CHAIN_ID="${CHAIN_ID:-irishub-2}"
NODE="${NODE:-tcp://127.0.0.1:26657}"
IRIS_HOME="${IRIS_HOME:-${HOME}/.iris}"
KEYRING_BACKEND="${KEYRING_BACKEND:-os}"
KEY_NAME="${KEY_NAME:-}"
UPGRADE_HEIGHT="${UPGRADE_HEIGHT:-}"
GAS_PRICES="${GAS_PRICES:-}"
UPGRADE_INFO="${UPGRADE_INFO:-}"
OUTPUT_DIR="${OUTPUT_DIR:-./irishub-v5-upgrade-proposal}"

fail() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

[[ -n "$KEY_NAME" ]] || fail "KEY_NAME is required"
[[ -n "$GAS_PRICES" ]] || fail "GAS_PRICES is required"
[[ "$UPGRADE_HEIGHT" =~ ^[1-9][0-9]*$ ]] || fail "UPGRADE_HEIGHT must be a positive integer"
if (( ${#UPGRADE_HEIGHT} > 19 )) ||
  { (( ${#UPGRADE_HEIGHT} == 19 )) && [[ "$UPGRADE_HEIGHT" > "9223372036854775807" ]]; }; then
  fail "UPGRADE_HEIGHT must fit in a signed 64-bit integer"
fi

command -v jq >/dev/null 2>&1 || fail "jq is required"
command -v "$IRIS_BIN" >/dev/null 2>&1 || fail "iris executable not found: $IRIS_BIN"
mkdir -p "$OUTPUT_DIR"

STATUS_FILE="$OUTPUT_DIR/node-status.json"
PARAMS_FILE="$OUTPUT_DIR/gov-params.json"
ACCOUNT_FILE="$OUTPUT_DIR/gov-account.json"
SUMMARY_FILE="$OUTPUT_DIR/irishub-v5-summary.md"
PROPOSAL_FILE="$OUTPUT_DIR/irishub-v5-proposal.json"
UNSIGNED_FILE="$OUTPUT_DIR/irishub-v5-unsigned-tx.json"
SUBMIT_OUTPUT_FILE="$OUTPUT_DIR/irishub-v5-submit-output.jsonl"
SUBMIT_FILE="$OUTPUT_DIR/irishub-v5-submit-response.json"

"$IRIS_BIN" status \
  --node "$NODE" \
  --home "$IRIS_HOME" \
  --output json > "$STATUS_FILE"

NODE_CHAIN_ID="$(
  jq -er '
    (.NodeInfo.network // .node_info.network)
    | select(type == "string" and length > 0)
  ' "$STATUS_FILE"
)" || fail "could not read the chain ID from $STATUS_FILE"
[[ "$NODE_CHAIN_ID" == "$CHAIN_ID" ]] ||
  fail "node chain ID is $NODE_CHAIN_ID, expected $CHAIN_ID"

LATEST_HEIGHT="$(
  jq -er '
    (.SyncInfo.latest_block_height // .sync_info.latest_block_height)
    | tostring
    | select(test("^[0-9]+$"))
  ' "$STATUS_FILE"
)" || fail "could not read the latest block height from $STATUS_FILE"
if (( ${#LATEST_HEIGHT} > ${#UPGRADE_HEIGHT} )) ||
  { (( ${#LATEST_HEIGHT} == ${#UPGRADE_HEIGHT} )) && [[ "$LATEST_HEIGHT" > "$UPGRADE_HEIGHT" || "$LATEST_HEIGHT" == "$UPGRADE_HEIGHT" ]]; }; then
  fail "UPGRADE_HEIGHT $UPGRADE_HEIGHT must be greater than latest block height $LATEST_HEIGHT"
fi

"$IRIS_BIN" query gov params \
  --node "$NODE" \
  --home "$IRIS_HOME" \
  --output json > "$PARAMS_FILE"

DEPOSIT="$(
  jq -er '
    .params.min_deposit
    | select(type == "array" and length > 0)
    | select(
        all(.[];
          (.amount | type == "string" and test("^[1-9][0-9]*$")) and
          (.denom | type == "string" and length > 0)
        )
      )
    | map(.amount + .denom)
    | join(",")
  ' "$PARAMS_FILE"
)" || fail "could not read a valid minimum deposit from $PARAMS_FILE"

"$IRIS_BIN" query auth module-account gov \
  --node "$NODE" \
  --home "$IRIS_HOME" \
  --output json > "$ACCOUNT_FILE"

GOV_AUTHORITY="$(
  jq -er '
    .account.base_account.address
    | select(type == "string" and length > 0)
  ' "$ACCOUNT_FILE"
)" || fail "could not read the governance module address from $ACCOUNT_FILE"

cat > "$SUMMARY_FILE" <<'EOF'
# Overview

This proposal schedules a coordinated upgrade of IRIS Hub to `v5.0.0`, using the on-chain upgrade plan name `v5.0.0`.

**If the proposal passes, validators will be required to update the IRIS Hub binary at the upgrade height specified in this proposal.**

**This upgrade is state-breaking.**

## Highlights

- Upgrade Cosmos SDK to `v0.53.8` and IBC-Go to `v10.7.0`.
- Remove the obsolete IBC capability store and its module-version entry.
- Execute the registered SDK and IBC module migrations, including the migration of IBC transfer denomination records to the new format while preserving existing IBC asset denominations.
- Retire legacy Liquid Staking Module (LSM) functionality and remove its obsolete staking records and fields.
- Preserve existing tokenized-share backing delegations, reward-owner associations, and derivative token balances and supplies.

### Breaking Changes

- Legacy LSM functionality is no longer available after the upgrade. Its records, indexes, and staking extensions are removed.
- Existing tokenized-share assets remain backed by their underlying delegations, but this release does not provide a transaction for redeeming those assets. A future redemption mechanism will require a separately agreed approach.
- Validators’ minimum self-delegation values are reset to zero to preserve the behavior of the previous LSM implementation.
- IBC integrations must support the interfaces and governance messages required by IBC-Go `v10.7.0`, including updated client recovery and parameter-update operations.

### Known Issues

- Redemption of existing LSM derivative tokens is not supported in this release. Their bank balances and supplies, underlying delegations, and reward-owner associations are preserved for future redemption handling.

## Release Details

- Target release: IRIS Hub `v5.0.0`.
- Upgrade plan name: `v5.0.0`.
- [Official IRIS Hub releases](https://github.com/irisnet/irishub/releases)
- [Changelog](https://github.com/irisnet/irishub/blob/master/CHANGELOG.md)

# Upgrade Instructions and Verification

Please check the announcements in the established IRIS Hub communication channels for the final release details, upgrade height, and validator instructions.

Validators must install the `v5.0.0` binary and restart their nodes at the scheduled upgrade height so that the registered upgrade handler can execute the required state migrations. Replacing the binary without executing the scheduled upgrade does not perform these migrations.

Before installation, verify the downloaded binary against the checksums and any verification instructions published with the official release.

# Voting

The voting period and proposal thresholds are governed by IRIS Hub’s active on-chain governance parameters.

The following items summarize the voting options and what they mean for this proposal:

YES - You agree that IRIS Hub should be upgraded to `v5.0.0` at the specified upgrade height.

NO - You disagree that IRIS Hub should be upgraded with this release.

NO WITH VETO - A `NoWithVeto` vote indicates that you consider the proposal to be spam, to disproportionately infringe on minority interests, or to violate the rules of engagement established by IRIS Hub governance. Veto thresholds and any deposit burning are governed by the chain’s active governance parameters.

ABSTAIN - You wish to contribute to the quorum but formally decline to vote either for or against the proposal.
EOF

jq -n \
  --arg authority "$GOV_AUTHORITY" \
  --arg height "$UPGRADE_HEIGHT" \
  --arg info "$UPGRADE_INFO" \
  --arg deposit "$DEPOSIT" \
  --rawfile summary "$SUMMARY_FILE" \
  '{
    messages: [
      {
        "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
        authority: $authority,
        plan: {
          name: "v5.0.0",
          height: $height,
          info: $info
        }
      }
    ],
    metadata: "",
    deposit: $deposit,
    title: "IRIS Hub v5.0.0 Software Upgrade",
    summary: $summary,
    expedited: false
  }' > "$PROPOSAL_FILE"

"$IRIS_BIN" tx gov submit-proposal "$PROPOSAL_FILE" \
  --from "$KEY_NAME" \
  --keyring-backend "$KEYRING_BACKEND" \
  --chain-id "$CHAIN_ID" \
  --node "$NODE" \
  --home "$IRIS_HOME" \
  --gas 500000 \
  --gas-prices "$GAS_PRICES" \
  --generate-only \
  --output json > "$UNSIGNED_FILE"

printf 'Prepared v5.0.0 proposal for %s at height %s.\n' "$CHAIN_ID" "$UPGRADE_HEIGHT"
printf 'Deposit: %s\nAuthority: %s\n' "$DEPOSIT" "$GOV_AUTHORITY"
printf 'Proposal: %s\nUnsigned transaction: %s\n' "$PROPOSAL_FILE" "$UNSIGNED_FILE"

if [[ "$submit" != true ]]; then
  printf 'No transaction was broadcast. Re-run with --submit after reviewing the files.\n'
  exit 0
fi

"$IRIS_BIN" tx gov submit-proposal "$PROPOSAL_FILE" \
  --from "$KEY_NAME" \
  --keyring-backend "$KEYRING_BACKEND" \
  --chain-id "$CHAIN_ID" \
  --node "$NODE" \
  --home "$IRIS_HOME" \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices "$GAS_PRICES" \
  --broadcast-mode sync \
  --output json | tee "$SUBMIT_OUTPUT_FILE"

jq -se '
  [
    .[]
    | select(
        (.txhash? | type == "string" and length > 0) and
        has("code")
      )
  ]
  | last // error("broadcast response not found")
' "$SUBMIT_OUTPUT_FILE" > "$SUBMIT_FILE" ||
  fail "broadcast output does not contain a transaction response; inspect $SUBMIT_OUTPUT_FILE"

TX_HASH="$(jq -er '.txhash | select(type == "string" and length > 0)' "$SUBMIT_FILE")" ||
  fail "broadcast response does not contain a transaction hash; inspect $SUBMIT_FILE"
BROADCAST_CODE="$(jq -er '.code | tonumber' "$SUBMIT_FILE")" ||
  fail "broadcast response contains an invalid result code; inspect $SUBMIT_FILE"
[[ "$BROADCAST_CODE" == 0 ]] ||
  fail "node rejected the transaction with code $BROADCAST_CODE; inspect $SUBMIT_FILE"

printf 'Broadcast response saved to %s. This does not confirm block inclusion.\n' "$SUBMIT_FILE"
printf 'After the transaction is included, verify it with:\n'
printf '  %q query tx %q --node %q --home %q --output json\n' \
  "$IRIS_BIN" "$TX_HASH" "$NODE" "$IRIS_HOME"
