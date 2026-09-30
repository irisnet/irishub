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
#   SIM_DIR         working directory (default: /tmp/irishub-sim). Needs
#                   roughly 3x the genesis size of free disk for the converted
#                   intermediates; do NOT point it at a tmpfs when the
#                   genesis is ~2GB
#   SIM_CHAIN_ID    simulation chain-id (default: irishub-sim-1); the final
#                   genesis carries this id because gentx signatures bind to it
#   KEEP_RUNNING=1  leave the node running when done (default: stop it)
#
# All large-file JSON conversions are streamed (O(1) memory, python3): the
# recovered genesis is ~2GB and jq would load 3-5x that into RAM. The node
# start itself is still memory-heavy for a 2GB genesis — on a 16GB machine run
# the script as:
#   GOMEMLIMIT=12GiB GOGC=30 bash scripts/chain-recovery/simulate-launch.sh ...
# (env vars are passed through to iris).
#
# Requires: bash, jq, python3, curl.

set -euo pipefail

IRIS_BIN=${IRIS_BIN:-iris}
EXPORT_GENESIS=${EXPORT_GENESIS:?set EXPORT_GENESIS to the recovered genesis file}
SIM_DIR=${SIM_DIR:-/tmp/irishub-sim}
SIM_CHAIN_ID=${SIM_CHAIN_ID:-irishub-sim-1}
KEEP_RUNNING=${KEEP_RUNNING:-0}
RPC=tcp://127.0.0.1:26657
RPC_HTTP=http://127.0.0.1:26657

die() { printf '%s\n' "$*" >&2; exit 1; }

# Streaming JSON member tool — trimmed copy of the one embedded in
# export-genesis.sh (same scanner core, verbatim), with one mode:
#   top <file> <key>                     print raw value of a top-level member
#   set-members <src> <dst> <member>=<value_file> ...
#                                       stream-copy src to dst, replacing the
#                                       listed members' values with the given
#                                       files' contents (complete JSON values)
#                                       and appending members missing from src
json_tool() {
  python3 - "$@" <<'PY'
import io, json, re, sys

CHUNK = 1 << 20
PAT = re.compile(rb'["\\\{\}\[\]]')
SCALAR_END = re.compile(rb'[,\}\] \t\n\r]')


class Stream:
    def __init__(self, path):
        self.f = open(path, 'rb')
        self.buf = b''
        self.pos = 0

    def fill(self):
        self.buf = self.f.read(CHUNK)
        self.pos = 0
        return len(self.buf)

    def peek(self):
        if self.pos >= len(self.buf) and self.fill() == 0:
            return b''
        return self.buf[self.pos:self.pos + 1]

    def adv(self):
        c = self.peek()
        if c:
            self.pos += 1
        return c

    def skipws(self):
        while self.peek() in (b' ', b'\t', b'\n', b'\r'):
            self.pos += 1


def scan_string(s, out):
    in_str = False
    # index (in the current chunk) of a byte consumed by a string escape, or -1
    esc_pos = -1
    rec = s.pos
    while True:
        buf = s.buf
        for m in PAT.finditer(buf, s.pos):
            if esc_pos >= 0 and m.start() == esc_pos:
                esc_pos = -1
                continue
            esc_pos = -1
            ch = m.group()
            if not in_str:
                if ch != b'"':
                    raise ValueError('bad string start %r' % ch)
                in_str = True
            elif ch == b'"':
                s.pos = m.end()
                if out is not None:
                    out.write(buf[rec:m.end()])
                return
            elif ch == b'\\':
                esc_pos = m.end()
        if out is not None:
            out.write(buf[rec:])
        if s.fill() == 0:
            raise ValueError('EOF inside string')
        # a backslash as the last byte of a chunk escapes the first byte of
        # the next chunk; otherwise the escaped byte (if any) is behind us
        esc_pos = 0 if esc_pos == len(buf) else -1
        rec = 0


def scan_container(s, out):
    depth = 0
    in_str = False
    # index (in the current chunk) of a byte consumed by a string escape, or -1
    esc_pos = -1
    rec = s.pos
    while True:
        buf = s.buf
        for m in PAT.finditer(buf, s.pos):
            if esc_pos >= 0 and m.start() == esc_pos:
                esc_pos = -1
                continue
            esc_pos = -1
            ch = m.group()
            if in_str:
                if ch == b'"':
                    in_str = False
                elif ch == b'\\':
                    esc_pos = m.end()
                continue
            if ch == b'"':
                in_str = True
            elif ch in (b'{', b'['):
                depth += 1
            else:
                depth -= 1
                if depth == 0:
                    s.pos = m.end()
                    if out is not None:
                        out.write(buf[rec:m.end()])
                    return
        if out is not None:
            out.write(buf[rec:])
        if s.fill() == 0:
            raise ValueError('EOF inside container')
        esc_pos = 0 if esc_pos == len(buf) else -1
        rec = 0


def scan_scalar(s, out):
    rec = s.pos
    any_char = False
    while True:
        buf = s.buf
        m = SCALAR_END.search(buf, s.pos)
        if m is not None:
            if not any_char and m.start() == s.pos:
                raise ValueError('empty scalar value')
            if out is not None:
                out.write(buf[rec:m.start()])
            s.pos = m.start()
            return
        if s.pos < len(buf):
            any_char = True
        if out is not None:
            out.write(buf[rec:])
        if s.fill() == 0:
            raise ValueError('EOF inside scalar')
        rec = 0


def scan_value(s, out=None):
    s.skipws()
    c = s.peek()
    if c == b'"':
        scan_string(s, out)
    elif c in (b'{', b'['):
        scan_container(s, out)
    elif c == b'':
        raise ValueError('EOF while expecting value')
    else:
        scan_scalar(s, out)


def read_key(s):
    s.skipws()
    if s.peek() != b'"':
        raise ValueError('expected object key')
    bio = io.BytesIO()
    scan_string(s, bio)
    return json.loads(bio.getvalue().decode('utf-8', 'surrogateescape'))


def iter_members(s):
    # yields keys; when a key is yielded, s is positioned at the start of its
    # value — the caller MUST consume the value with scan_value(s, ...)
    s.skipws()
    if s.adv() != b'{':
        raise ValueError('expected object')
    s.skipws()
    if s.peek() == b'}':
        s.adv()
        return
    while True:
        key = read_key(s)
        s.skipws()
        if s.adv() != b':':
            raise ValueError('expected colon')
        s.skipws()
        yield key
        s.skipws()
        c = s.adv()
        if c == b',':
            continue
        if c == b'}':
            return
        raise ValueError('bad object structure')


def cmd_top(path, key):
    s = Stream(path)
    for k in iter_members(s):
        if k == key:
            scan_value(s, sys.stdout.buffer)
            sys.stdout.buffer.write(b'\n')
            return 0
        scan_value(s)
    sys.stderr.write('member not found: %s\n' % key)
    return 3


def cmd_set_members(src, dst, pairs):
    # stream-copy src to dst, replacing the listed top-level members' values
    # with the contents of the given files (complete JSON values); members
    # not present in src are appended at the end of the object
    repl = {}
    for p in pairs:
        member, sep, path = p.partition('=')
        if not member or not sep or not path:
            raise ValueError('expected member=value_file, got %r' % p)
        repl[member] = path
    afs = {m: open(f, 'rb') for m, f in repl.items()}
    s = Stream(src)
    out = open(dst, 'wb')
    try:
        out.write(b'{')
        wrote = False
        for k in iter_members(s):
            if wrote:
                out.write(b',\n')
            wrote = True
            out.write(json.dumps(k).encode('utf-8'))
            out.write(b':')
            if k in repl:
                af = afs.pop(k)
                try:
                    while True:
                        b = af.read(CHUNK)
                        if not b:
                            break
                        out.write(b)
                finally:
                    af.close()
                scan_value(s, None)
            else:
                scan_value(s, out)
        for m, af in afs.items():
            if wrote:
                out.write(b',\n')
            wrote = True
            out.write(json.dumps(m).encode('utf-8'))
            out.write(b':')
            try:
                while True:
                    b = af.read(CHUNK)
                    if not b:
                        break
                    out.write(b)
            finally:
                af.close()
        out.write(b'}\n')
    finally:
        out.close()
        for af in afs.values():
            af.close()
    return 0


def main():
    if len(sys.argv) < 4:
        sys.stderr.write('usage: json_tool MODE ARGS...\n')
        return 2
    mode = sys.argv[1]
    try:
        if mode == 'top' and len(sys.argv) == 4:
            return cmd_top(sys.argv[2], sys.argv[3])
        if mode == 'set-members' and len(sys.argv) >= 5:
            return cmd_set_members(sys.argv[2], sys.argv[3], sys.argv[4:])
    except (ValueError, KeyError) as e:
        sys.stderr.write('json_tool %s: %s\n' % (mode, e))
        return 1
    sys.stderr.write('bad arguments for mode %s\n' % mode)
    return 2


sys.exit(main())
PY
}

# to_comet <src> <dst>: rewrite an SDK genesis in cometbft's strict JSON
# spelling — `replace-validators` reads genesis files with cometbft's decoder,
# which rejects the SDK's numeric initial_height and null app_hash. Single
# streaming pass, O(1) memory.
to_comet() {
  local src=$1 dst=$2 ih ah
  ih=$(json_tool top "$src" initial_height) || return 1
  ih=${ih#\"}; ih=${ih%\"}
  ah=$(json_tool top "$src" app_hash) || return 1
  if [[ $ah == null ]]; then ah='""'; fi
  printf '"%s"' "$ih" > "$dst.ih.tmp"
  printf '%s' "$ah" > "$dst.ah.tmp"
  json_tool set-members "$src" "$dst" \
    initial_height="$dst.ih.tmp" app_hash="$dst.ah.tmp" || return 1
  rm -f "$dst.ih.tmp" "$dst.ah.tmp"
}

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

printf '==> converting the genesis files to cometbft spelling\n'
to_comet "$NODE_HOME/config/genesis.json" "$SIM_DIR/source.comet.json" \
  || die "failed to convert the testnet genesis"
to_comet "$EXPORT_GENESIS" "$SIM_DIR/target.comet.json" \
  || die "failed to convert the recovered genesis"

printf '==> replacing the exported validator set with ours\n'
"$IRIS_BIN" genesis replace-validators \
  --source-genesis-file "$SIM_DIR/source.comet.json" \
  --target-genesis-file "$SIM_DIR/target.comet.json" \
  --output-genesis-file "$SIM_DIR/simulation.comet.json"

printf '==> converting the merged genesis back to SDK format\n'
# re-attach the SDK-only members (app_name/app_version/consensus) the comet
# decoder dropped; consensus.validators MUST stay empty: the doc validator set
# would otherwise mix the old (uncontrolled) validators with ours and block
# consensus, the set comes from the imported gentxs via InitChain instead
an=$(json_tool top "$EXPORT_GENESIS" app_name) \
  || die "no app_name in $EXPORT_GENESIS"
av=$(json_tool top "$EXPORT_GENESIS" app_version) \
  || die "no app_version in $EXPORT_GENESIS"
printf '%s' "$an" > "$SIM_DIR/app-name.json"
printf '%s' "$av" > "$SIM_DIR/app-version.json"
json_tool top "$EXPORT_GENESIS" consensus > "$SIM_DIR/consensus.json" \
  || die "no consensus in $EXPORT_GENESIS"
python3 - "$SIM_DIR/consensus.json" "$SIM_DIR/consensus-final.json" <<'PY2'
import json, sys
cons = json.load(open(sys.argv[1]))
json.dump({'params': cons.get('params'), 'validators': []},
          open(sys.argv[2], 'w'))
PY2
ih=$(json_tool top "$SIM_DIR/simulation.comet.json" initial_height) \
  || die "no initial_height in the merged genesis"
ih=${ih#\"}; ih=${ih%\"}   # cometbft writes int64s quoted
printf '%s' "$ih" > "$SIM_DIR/initial-height.json"
json_tool set-members "$SIM_DIR/simulation.comet.json" \
  "$SIM_DIR/simulation-genesis.json" \
  initial_height="$SIM_DIR/initial-height.json" \
  app_name="$SIM_DIR/app-name.json" \
  app_version="$SIM_DIR/app-version.json" \
  consensus="$SIM_DIR/consensus-final.json" \
  || die "failed to convert the merged genesis back to SDK format"
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
