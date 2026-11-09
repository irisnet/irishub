#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
script="$script_dir/submit-v5-upgrade-proposal.sh"
test_root="$(mktemp -d /tmp/irishub-v5-proposal-test.XXXXXX)"
trap 'rm -rf "$test_root"' EXIT
mkdir -p "$test_root/bin"

cat > "$test_root/bin/iris" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf '%q ' "$@" >> "$MOCK_LOG"
printf '\n' >> "$MOCK_LOG"
case "${1:-}" in
  status)
    [[ "${MOCK_STATUS_FAIL:-false}" != true ]] || exit 1
    printf '{"NodeInfo":{"network":"irishub-2"},"SyncInfo":{"latest_block_height":"100"}}\n'
    ;;
  query)
    case "${2:-} ${3:-} ${4:-}" in
      "gov params --node")
        [[ "${MOCK_PARAMS_FAIL:-false}" != true ]] || exit 1
        if [[ "${MOCK_PARAMS_MALFORMED:-false}" == true ]]; then
          printf 'warning: malformed response\n{"params":' >&2
          exit 0
        fi
        printf '{"params":{"min_deposit":[{"denom":"unyan","amount":"10000000"}]}}\n' >&2
        ;;
      "auth module-account gov")
        printf '{"account":{"base_account":{"address":"iaa1govauthority"}}}\n' >&2
        ;;
      *) exit 3 ;;
    esac
    ;;
  tx)
    if printf '%s\n' "$@" | grep -qx -- '--generate-only'; then
      printf '{"body":{"messages":[{"@type":"/cosmos.gov.v1.MsgSubmitProposal"}]}}\n'
    else
      printf '{"body":{"messages":[{"@type":"/cosmos.gov.v1.MsgSubmitProposal"}]}}\n'
      printf 'gas estimate: 321000\n' >&2
      printf 'confirm transaction before signing and broadcasting [y/N]: ' >&2
      [[ "${MOCK_CANCEL:-false}" != true ]] || exit 1
      printf '{"txhash":"ABC123","code":0}\n'
    fi
    ;;
  *) exit 4 ;;
esac
EOF
chmod +x "$test_root/bin/iris"

run_script() {
  env PATH="$test_root/bin:$PATH" \
    MOCK_LOG="$1" \
    IRIS_BIN=iris \
    KEY_NAME=proposer \
    UPGRADE_HEIGHT="${2:-200}" \
    GAS_PRICES=0.3uiris \
    OUTPUT_DIR="$3" \
    "$script" "${@:4}"
}

default_log="$test_root/default.log"
default_out="$test_root/default"
run_script "$default_log" 200 "$default_out"
[[ "$(grep -c '^tx gov submit-proposal' "$default_log")" == 1 ]]
grep -q -- '--generate-only' "$default_log"
[[ ! -e "$default_out/irishub-v5-submit-response.json" ]]
jq -e '
  .messages[0]["@type"] == "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade" and
  .messages[0].authority == "iaa1govauthority" and
  .messages[0].plan.name == "v5.0.0" and
  .messages[0].plan.height == "200" and
  .messages[0].plan.info == "" and
  .deposit == "10000000unyan" and
  .title == "IRIS Hub v5.0.0 Software Upgrade" and
  .expedited == false and
  (.summary | startswith("# Overview\n")) and
  (.summary | contains("NO WITH VETO"))
' "$default_out/irishub-v5-proposal.json" >/dev/null
cmp -s "$default_out/irishub-v5-summary.md" \
  <(jq -rj '.summary' "$default_out/irishub-v5-proposal.json")

submit_log="$test_root/submit.log"
submit_out="$test_root/submit"
submit_stderr="$test_root/submit.stderr"
run_script "$submit_log" 200 "$submit_out" --submit 2> "$submit_stderr"
[[ "$(grep -c '^tx gov submit-proposal' "$submit_log")" == 2 ]]
[[ "$(grep -c -- '--generate-only' "$submit_log")" == 1 ]]
! grep -Eq -- '(^| )(-y|--yes)( |$)' "$submit_log"
jq -e '.txhash == "ABC123" and .code == 0' \
  "$submit_out/irishub-v5-submit-response.json" >/dev/null
[[ "$(jq -s 'length' "$submit_out/irishub-v5-submit-output.jsonl")" == 2 ]]
grep -q 'gas estimate: 321000' "$submit_stderr"

cancel_log="$test_root/cancel.log"
cancel_out="$test_root/cancel"
cancel_stderr="$test_root/cancel.stderr"
if env PATH="$test_root/bin:$PATH" MOCK_LOG="$cancel_log" MOCK_CANCEL=true \
  IRIS_BIN=iris KEY_NAME=proposer UPGRADE_HEIGHT=200 GAS_PRICES=0.3uiris \
  OUTPUT_DIR="$cancel_out" "$script" --submit >/dev/null 2> "$cancel_stderr"; then
  printf 'expected canceled submission to fail\n' >&2
  exit 1
fi
[[ ! -e "$cancel_out/irishub-v5-submit-response.json" ]]
[[ "$(jq -s 'length' "$cancel_out/irishub-v5-submit-output.jsonl")" == 1 ]]
grep -q 'confirm transaction' "$cancel_stderr"

invalid_log="$test_root/invalid.log"
if run_script "$invalid_log" invalid "$test_root/invalid" >/dev/null 2>&1; then
  printf 'expected invalid height to fail\n' >&2
  exit 1
fi
[[ ! -s "$invalid_log" ]]

stale_log="$test_root/stale.log"
if run_script "$stale_log" 100 "$test_root/stale" >/dev/null 2>&1; then
  printf 'expected non-future height to fail\n' >&2
  exit 1
fi
[[ "$(wc -l < "$stale_log" | tr -d ' ')" == 1 ]]

failed_log="$test_root/failed.log"
if env PATH="$test_root/bin:$PATH" MOCK_LOG="$failed_log" MOCK_PARAMS_MALFORMED=true \
  IRIS_BIN=iris KEY_NAME=proposer UPGRADE_HEIGHT=200 GAS_PRICES=0.3uiris \
  OUTPUT_DIR="$test_root/failed" "$script" >/dev/null 2>&1; then
  printf 'expected malformed governance parameters to stop the script\n' >&2
  exit 1
fi
[[ "$(wc -l < "$failed_log" | tr -d ' ')" == 2 ]]
! grep -q '^tx ' "$failed_log"

printf 'submit-v5-upgrade-proposal tests passed\n'
