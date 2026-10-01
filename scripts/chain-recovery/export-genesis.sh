#!/usr/bin/env bash
#
# Memory-bounded genesis export for very large chains.
#
# The stock `iris export` holds ~3-5x the full app-state JSON in RAM at once
# (the module manager keeps every module's genesis in one map, then
# json.MarshalIndent copies it, then json.Marshal(appGenesis) copies it again),
# which OOMs on large data — swap only turns it into a GC livelock.
#
# This script instead:
#   1. runs `iris export` once per module (default) with --modules-to-export,
#      so peak RSS is bounded by the largest single module; resumes are
#      content-addressed (partials are matched to groups by their module set,
#      so the module list may change between runs without invalidating them);
#   2. assembles the final genesis with a streaming JSON member scanner
#      (python, O(1) memory), producing the same result as the old jq merge:
#      chain_id replaced, cross-chain modules (IBC + TIBC) reset to
#      fresh-init defaults, everything else taken from the export at
#      EXPORT_HEIGHT; modules that export no genesis (skipped by the module
#      manager because they implement no genesis interfaces) fall back to
#      their fresh default;
#   3. sanitizes legacy chain data that the v0.50 genesis validation rejects
#      (blank bank denom_metadata name/symbol, unusable denom units, gov
#      expedited_min_deposit <= min_deposit), logging every change;
#   4. validates each exported module by building a "padded" genesis — the
#      fresh defaults with that one module replaced by the exported fragment
#      (`iris genesis validate` requires every module section to be present,
#      so a bare partial cannot be validated directly) — and reports ALL
#      module failures at once;
#   5. asserts chain_id / initial_height / uiris supply with streaming tools
#      (no jq slurpfile on the huge final file).
#
# Note: `iris export` runs service.PrepForZeroHeightGenesis on every run
# (app/export.go), so it is paid once per group here.
#
# Requires: bash, jq, python3.
#
# Env:
#   IRIS_BIN, OLD_HOME, OUT_DIR,
#   NEW_CHAIN_ID, EXPORT_HEIGHT, EXPECTED_INITIAL_HEIGHT, EXPECTED_SUPPLY
#       — same meaning as before.
#   MODULE_GROUPS  optional: newline-separated groups of comma-joined module
#                 names, one `iris export` run per group. Default: one module
#                 per run (safest memory bound). Use it to merge small modules
#                 into fewer runs once partial sizes are known.
#   RESUME=1      reuse existing partials from an interrupted run (assumes the
#                 same MODULE_GROUPS grouping).
#   SKIP_VALIDATE=1  skip `iris genesis validate` calls (assertions still run).

set -euo pipefail

IRIS_BIN=${IRIS_BIN:-/root/iris}
OLD_HOME=${OLD_HOME:-/data/iris}
OUT_DIR=${OUT_DIR:-/data/export-output}

NEW_CHAIN_ID=${NEW_CHAIN_ID:-irishub-2}
EXPORT_HEIGHT=${EXPORT_HEIGHT:-37242247}
EXPECTED_INITIAL_HEIGHT=${EXPECTED_INITIAL_HEIGHT:-37242248}
EXPECTED_SUPPLY=${EXPECTED_SUPPLY:-2156713998266827}

# cross-chain modules whose state is reset to fresh-init defaults instead of
# being exported: the IBC/TIBC client & connection state of the legacy chain
# is not carried over to the recovered chain
RESET_MODULES=(07-tendermint ibc transfer interchainaccounts nonfungibletokentransfer tibc)
RESUME=${RESUME:-0}
SKIP_VALIDATE=${SKIP_VALIDATE:-0}
MODULE_GROUPS=${MODULE_GROUPS:-}

WORK="$OUT_DIR/work"
FRAG="$WORK/frag"
FRESH="$OUT_DIR/config/genesis.json"

die() { printf '%s\n' "$*" >&2; exit 1; }

# Streaming JSON member tool: parses huge JSON without loading it.
#   top <file> <key>                     print raw value of a top-level member
#   split <file> <container> <outdir>    write each member of <container>'s
#                                        value to <outdir>/<member>.json and
#                                        print the member names
#   patch <src> <dst> <member> <value_file>
#                                       stream-copy <src> to <dst>, replacing
#                                       the top-level member's value with the
#                                       contents of <value_file>
#   skeleton <src> <dst> <chain_id> <app_state_file>
#                                       copy <src> verbatim to <dst> but with
#                                       chain_id replaced and app_state value
#                                       replaced by <app_state_file> contents
#   supply <file> <denom>                print amount of <denom> in
#                                        app_state.bank.supply
json_tool() {
  python3 - "$@" <<'PY'
import io, json, os, re, sys

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


def iter_array(s):
    # yields once per item; s is positioned at the item value — the caller
    # MUST consume it with scan_value(s, ...)
    s.skipws()
    if s.adv() != b'[':
        raise ValueError('expected array')
    s.skipws()
    if s.peek() == b']':
        s.adv()
        return
    while True:
        s.skipws()
        yield
        s.skipws()
        c = s.adv()
        if c == b',':
            continue
        if c == b']':
            return
        raise ValueError('bad array structure')


def find_path(s, path):
    for i, key in enumerate(path):
        found = False
        for k in iter_members(s):
            if k == key:
                found = True
                break
            scan_value(s)
        if not found:
            raise KeyError('member not found: ' + '.'.join(path[:i + 1]))


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


def cmd_split(path, container, outdir):
    s = Stream(path)
    for k in iter_members(s):
        if k == container:
            # fragment files are counter-named (not named after the module)
            # because module names may differ only by case, which collides
            # on case-insensitive filesystems; the counter continues across
            # invocations sharing the directory; the mapping is printed as
            # "<module>\t<file>" lines on stdout
            idx = len(os.listdir(outdir))
            out = []
            for k2 in iter_members(s):
                idx += 1
                fn = '%05d.json' % idx
                with open(os.path.join(outdir, fn), 'wb') as o:
                    scan_value(s, o)
                out.append(k2 + '\t' + fn)
            for line in out:
                sys.stdout.write(line + '\n')
            return 0
        scan_value(s)
    sys.stderr.write('container not found: %s\n' % container)
    return 3


def cmd_subkeys(path, container):
    s = Stream(path)
    for k in iter_members(s):
        if k == container:
            for k2 in iter_members(s):
                sys.stdout.write(k2 + '\n')
                scan_value(s)
            return 0
        scan_value(s)
    sys.stderr.write('container not found: %s\n' % container)
    return 3


CHARSET = 'qpzry9x8gf2tvdw0s3jn54khce6mua7l'


def bech32_to_bytes(addr):
    # decode the 20-byte payload of an iaa1... address (checksum not verified)
    try:
        pos = addr.index('1')
        acc = 0
        bits = 0
        out = bytearray()
        for c in addr[pos + 1:]:
            v = CHARSET.find(c)
            if v < 0:
                return None
            acc = (acc << 5) | v
            bits += 5
            while bits >= 8:
                bits -= 8
                out.append((acc >> bits) & 0xff)
        return bytes(out[:20]) if len(out) >= 20 else None
    except ValueError:
        return None


def collect_auth_addrs(s, keep):
    # collect the hex form of every "address" inside an auth genesis value
    for m in iter_members(s):
        if m == 'accounts':
            for _ in iter_array(s):
                bio = io.BytesIO()
                scan_value(s, bio)
                am = re.search(rb'"address":\s*"([^"]+)"', bio.getvalue())
                if am:
                    raw = bech32_to_bytes(am.group(1).decode())
                    if raw:
                        keep.add(raw.hex())
        else:
            scan_value(s)


def cmd_filter_evm(src, dst):
    # drop evm genesis accounts whose address has no matching auth account:
    # such entries carry unreachable code/storage and evm InitGenesis panics
    # on the missing account
    keep = set()
    s = Stream(src)
    for k in iter_members(s):
        if k == 'app_state':
            for m in iter_members(s):
                if m == 'auth':
                    collect_auth_addrs(s, keep)
                else:
                    scan_value(s)
            break
        scan_value(s)

    dropped = 0
    total = 0

    def rewrite_evm(s, out):
        nonlocal dropped, total
        out.write(b'{')
        wrote = False
        for m in iter_members(s):
            if wrote:
                out.write(b',')
            wrote = True
            out.write(json.dumps(m).encode('utf-8'))
            out.write(b':')
            if m == 'accounts':
                out.write(b'[')
                first = True
                for _ in iter_array(s):
                    bio = io.BytesIO()
                    scan_value(s, bio)
                    el = bio.getvalue()
                    total += 1
                    am = re.search(rb'"address":\s*"0x([0-9a-fA-F]{40})"', el)
                    key = am.group(1).decode().lower() if am else None
                    if key is None or key not in keep:
                        dropped += 1
                        print('evm: dropped account %s (no matching auth account)'
                              % (key or '<no address>'), file=sys.stderr)
                        continue
                    if not first:
                        out.write(b',')
                    first = False
                    out.write(el)
                out.write(b']')
            else:
                scan_value(s, out)
        out.write(b'}')

    s = Stream(src)
    out = open(dst, 'wb')
    out.write(b'{')
    wrote = False
    with out:
        for k in iter_members(s):
            if wrote:
                out.write(b',\n')
            wrote = True
            out.write(json.dumps(k).encode('utf-8'))
            out.write(b':')
            if k == 'app_state':
                out.write(b'{')
                wrote2 = False
                for m in iter_members(s):
                    if wrote2:
                        out.write(b',')
                    wrote2 = True
                    out.write(json.dumps(m).encode('utf-8'))
                    out.write(b':')
                    if m == 'evm':
                        rewrite_evm(s, out)
                    else:
                        scan_value(s, out)
                out.write(b'}')
            else:
                scan_value(s, out)
        out.write(b'}\n')

    print('evm sanitize: %d auth account(s), %d evm account(s), %d dropped'
          % (len(keep), total, dropped), file=sys.stderr)
    return 0


def cmd_fix_liquid_stake(src, dst):
    # normalize the staking genesis' total_liquid_staked_tokens counter to the
    # sum of the validators' liquid_shares: legacy chains can carry a 1-unit
    # drift between the cached counter and the per-validator state, which the
    # crisis module's liquid stake invariant panics on at genesis import
    s = Stream(src)
    for k in iter_members(s):
        if k == 'app_state':
            for m in iter_members(s):
                if m == 'staking':
                    return fix_staking_member(s, dst)
                scan_value(s)
            break
        scan_value(s)
    sys.stderr.write('staking member not found\n')
    return 3


def fix_staking_member(s, dst):
    # single pass: sum liquid_shares while streaming the validators array,
    # rewrite total_liquid_staked_tokens when reached
    total = [None]  # mutable box for the nested writer below

    def rewrite_staking(s, out):
        out.write(b'{')
        wrote = False
        for m in iter_members(s):
            if wrote:
                out.write(b',')
            wrote = True
            out.write(json.dumps(m).encode('utf-8'))
            out.write(b':')
            if m == 'validators':
                out.write(b'[')
                first = True
                for _ in iter_array(s):
                    bio = io.BytesIO()
                    scan_value(s, bio)
                    el = bio.getvalue()
                    # replicate the invariant: sum over all validators of
                    # int(liquid_shares * tokens / delegator_shares) with
                    # sdk.Dec semantics (18-decimal truncating arithmetic)
                    ls = re.search(rb'"liquid_shares":\s*"([0-9.]+)"', el)
                    if ls:
                        from decimal import Decimal, ROUND_DOWN
                        if total[0] is None:
                            total[0] = 0
                        tk = re.search(rb'"tokens":\s*"([0-9.]+)"', el)
                        dsh = re.search(rb'"delegator_shares":\s*"([0-9.]+)"', el)
                        ls_d = Decimal(ls.group(1).decode())
                        dsh_d = Decimal(dsh.group(1).decode()) if dsh else Decimal(0)
                        tk_d = Decimal(tk.group(1).decode()) if tk else Decimal(0)
                        if dsh_d != 0:
                            q = (ls_d / dsh_d).quantize(Decimal('1e-18'), rounding=ROUND_DOWN)
                            contrib = int((q * tk_d).quantize(Decimal('1e-18'), rounding=ROUND_DOWN))
                            total[0] += contrib
                    if not first:
                        out.write(b',')
                    first = False
                    out.write(el)
                out.write(b']')
            elif m == 'total_liquid_staked_tokens':
                bio = io.BytesIO()
                scan_value(s, bio)
                old = bio.getvalue()
                if total[0] is None:
                    out.write(old)
                else:
                    new = ('"%d"' % int(total[0])).encode()
                    if old != new:
                        print('staking: total_liquid_staked_tokens %s -> %s '
                              '(sum of validators\' liquid_shares)'
                              % (old.decode(), new.decode().strip('"')), file=sys.stderr)
                    out.write(new)
            else:
                scan_value(s, out)
        out.write(b'}')

    s2 = Stream(s.f.name)
    out = open(dst, 'wb')
    out.write(b'{')
    wrote = False
    with out:
        for k in iter_members(s2):
            if wrote:
                out.write(b',\n')
            wrote = True
            out.write(json.dumps(k).encode('utf-8'))
            out.write(b':')
            if k == 'app_state':
                out.write(b'{')
                wrote2 = False
                for m in iter_members(s2):
                    if wrote2:
                        out.write(b',')
                    wrote2 = True
                    out.write(json.dumps(m).encode('utf-8'))
                    out.write(b':')
                    if m == 'staking':
                        rewrite_staking(s2, out)
                    else:
                        scan_value(s2, out)
                out.write(b'}')
            else:
                scan_value(s2, out)
        out.write(b'}\n')
    if total[0] is None:
        sys.stderr.write('staking: no validators seen, nothing to normalize\n')
    return 0


def cmd_patch(src, dst, member, value_file):
    # stream-copy src to dst, replacing the top-level member's value with the
    # contents of value_file (a complete JSON value)
    s = Stream(src)
    out = open(dst, 'wb')
    out.write(b'{')
    wrote = False
    with open(value_file, 'rb') as af, out:
        for k in iter_members(s):
            if wrote:
                out.write(b',\n')
            out.write(json.dumps(k).encode('utf-8'))
            out.write(b':')
            if k == member:
                while True:
                    b = af.read(CHUNK)
                    if not b:
                        break
                    out.write(b)
                scan_value(s, None)
            else:
                scan_value(s, out)
            wrote = True
        out.write(b'}\n')
    return 0


def cmd_skeleton(src, dst, chain_id, app_state_file):
    s = Stream(src)
    out = open(dst, 'wb')
    out.write(b'{')
    wrote = False
    with open(app_state_file, 'rb') as af, out:
        for k in iter_members(s):
            if wrote:
                out.write(b',\n')
            out.write(json.dumps(k).encode('utf-8'))
            out.write(b':')
            if k == 'chain_id':
                out.write(json.dumps(chain_id).encode('utf-8'))
                scan_value(s, None)
            elif k == 'app_state':
                while True:
                    b = af.read(CHUNK)
                    if not b:
                        break
                    out.write(b)
                scan_value(s, None)
            else:
                scan_value(s, out)
            wrote = True
        out.write(b'}\n')
    return 0


def cmd_supply(path, denom):
    s = Stream(path)
    try:
        find_path(s, ['app_state', 'bank', 'supply'])
    except KeyError as e:
        sys.stderr.write('%s\n' % e)
        return 3
    for _ in iter_array(s):
        bio = io.BytesIO()
        scan_value(s, bio)
        item = json.loads(bio.getvalue().decode('utf-8', 'surrogateescape'))
        if isinstance(item, dict) and item.get('denom') == denom:
            amt = item.get('amount')
            if amt is None:
                sys.stderr.write('supply entry for %s has no amount\n' % denom)
                return 1
            print(amt)
            return 0
    sys.stderr.write('denom not found in supply: %s\n' % denom)
    return 4


def main():
    if len(sys.argv) < 3:
        sys.stderr.write('usage: json_tool MODE ARGS...\n')
        return 2
    mode = sys.argv[1]
    try:
        if mode == 'top' and len(sys.argv) == 4:
            return cmd_top(sys.argv[2], sys.argv[3])
        if mode == 'split' and len(sys.argv) == 5:
            return cmd_split(sys.argv[2], sys.argv[3], sys.argv[4])
        if mode == 'subkeys' and len(sys.argv) == 4:
            return cmd_subkeys(sys.argv[2], sys.argv[3])
        if mode == 'filter-evm' and len(sys.argv) == 4:
            return cmd_filter_evm(sys.argv[2], sys.argv[3])
        if mode == 'fix-liquid-stake' and len(sys.argv) == 4:
            return cmd_fix_liquid_stake(sys.argv[2], sys.argv[3])
        if mode == 'skeleton' and len(sys.argv) == 6:
            return cmd_skeleton(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5])
        if mode == 'patch' and len(sys.argv) == 6:
            return cmd_patch(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5])
        if mode == 'supply' and len(sys.argv) == 4:
            return cmd_supply(sys.argv[2], sys.argv[3])
    except (ValueError, KeyError) as e:
        sys.stderr.write('json_tool %s: %s\n' % (mode, e))
        return 1
    sys.stderr.write('bad arguments for mode %s\n' % mode)
    return 2


sys.exit(main())
PY
}

### guards & init #############################################################

if [[ $RESUME == 1 ]]; then
  if [[ -e $OUT_DIR/genesis.json ]]; then
    die "refusing to overwrite existing path: $OUT_DIR/genesis.json"
  fi
  if [[ ! -s $FRESH ]]; then
    "$IRIS_BIN" init recovery \
      --chain-id "$NEW_CHAIN_ID" \
      --home "$OUT_DIR" >/dev/null
  fi
else
  for output in "$OUT_DIR/config/genesis.json" "$OUT_DIR/partial.json" \
    "$OUT_DIR/genesis.json" "$WORK"; do
    if [[ -e "$output" ]]; then
      printf 'refusing to overwrite existing path: %s\n' "$output" >&2
      exit 1
    fi
  done
  mkdir -p "$OUT_DIR"
  "$IRIS_BIN" init recovery \
    --chain-id "$NEW_CHAIN_ID" \
    --home "$OUT_DIR" >/dev/null
fi
mkdir -p "$WORK" "$FRAG"

### module list (export everything except the cross-chain modules, which are
### reset to fresh-init defaults in the final genesis) ########################

mapfile -t ALL_KEYS < <(jq -r '.app_state | keys[]' "$FRESH" | sort)
declare -A IBC_SET=()
for m in "${RESET_MODULES[@]}"; do IBC_SET[$m]=1; done

MODULES=()
for m in "${ALL_KEYS[@]}"; do
  if [[ -z ${IBC_SET[$m]:-} ]]; then
    MODULES+=("$m")
  fi
done
((${#MODULES[@]} > 0)) || die "no modules to export"

for m in "${MODULES[@]}" "${RESET_MODULES[@]}"; do
  [[ $m =~ ^[A-Za-z0-9._-]+$ ]] || die "unexpected module name: $m"
done

### batched export ###########################################################

mapfile -t EXPORT_GROUPS < <(
  if [[ -n $MODULE_GROUPS ]]; then
    printf '%s\n' "$MODULE_GROUPS"
  else
    printf '%s\n' "${MODULES[@]}"
  fi
)

# content-addressed resume: match existing partials to groups by their
# app_state module set instead of by index — the module list can change
# between runs (e.g. a module moved to the reset list), which would shift
# positional numbering and silently reuse the wrong partial. Partials that
# belong to no group (e.g. a dropped module from an earlier run) are pruned.
declare -A PARTIAL_BY_KEY=()
if [[ $RESUME == 1 ]]; then
  for p in "$WORK"/partial-*.json "$WORK"/stage*/partial-*.json; do
    [[ -e $p ]] || continue
    key=$(json_tool subkeys "$p" app_state | sort | tr '\n' ',')
    # empty key = a module that exports no genesis; it cannot be matched by
    # content and is re-exported (cheap)
    [[ -n $key ]] || continue
    PARTIAL_BY_KEY[$key]=$p
  done
fi

STAGE="$WORK/stage.$$"
mkdir -p "$STAGE"

i=0
for g in "${EXPORT_GROUPS[@]}"; do
  [[ -n $g ]] || continue
  i=$((i+1))
  key=$(printf '%s\n' ${g//,/ } | sort | tr '\n' ',')
  if [[ -n ${PARTIAL_BY_KEY[$key]:-} ]]; then
    printf 'resume: reusing %s for group %d (%s)\n' \
      "${PARTIAL_BY_KEY[$key]}" "$i" "$g"
    mv "${PARTIAL_BY_KEY[$key]}" "$STAGE/partial-$i.json"
  else
    printf 'export group %d: %s\n' "$i" "$g"
    "$IRIS_BIN" export \
      --home "$OLD_HOME" \
      --height "$EXPORT_HEIGHT" \
      --for-zero-height=false \
      --modules-to-export "$g" \
      --output-document "$STAGE/partial-$i.json"
  fi

  ih=$(json_tool top "$STAGE/partial-$i.json" initial_height) \
    || die "failed to read initial_height from partial-$i.json"
  [[ $ih == "$EXPECTED_INITIAL_HEIGHT" ]] \
    || die "partial-$i.json: initial_height=$ih, want $EXPECTED_INITIAL_HEIGHT"
done

# swap in this run's partials; whatever was left in the pool or in old
# interrupted staging dirs belongs to no group (e.g. a module dropped from
# the module list since an earlier run) and is pruned
rm -f "$WORK"/partial-*.json
mv "$STAGE"/partial-*.json "$WORK/" || die "failed to move partials back"
rm -rf "$WORK"/stage*

### streaming assembly #######################################################

rm -rf "$FRAG"
mkdir -p "$FRAG"
declare -A FRAG_FILE=()
: > "$WORK/exported.txt"
for partial in "$WORK"/partial-*.json; do
  [[ -e $partial ]] || continue
  json_tool split "$partial" app_state "$FRAG" >> "$WORK/exported.txt" \
    || die "failed to split $partial"
done
while IFS=$'\t' read -r mod fn; do
  FRAG_FILE[$mod]="$FRAG/$fn"
done < "$WORK/exported.txt"
mapfile -t EXPORTED < <(cut -f1 "$WORK/exported.txt" | sort -u)

# fresh defaults for every module: used to reset the IBC modules in the final
# genesis, to pad partials for `iris genesis validate`, and as a fallback for
# modules that export no genesis
DEF="$WORK/defaults"
rm -rf "$DEF"
mkdir -p "$DEF"
declare -A DEF_FILE=()
json_tool split "$FRESH" app_state "$DEF" > "$WORK/defaults.txt" \
  || die "failed to split fresh genesis defaults"
while IFS=$'\t' read -r mod fn; do
  DEF_FILE[$mod]="$DEF/$fn"
done < "$WORK/defaults.txt"
for m in "${ALL_KEYS[@]}"; do
  [[ -n ${DEF_FILE[$m]:-} ]] || die "fresh defaults do not cover module $m"
done

# nothing unexpected may appear in the export output
declare -A MOD_SET=()
for m in "${MODULES[@]}"; do MOD_SET[$m]=1; done
for m in "${EXPORTED[@]}"; do
  [[ -n ${MOD_SET[$m]:-} ]] || die "unexpected exported module: $m"
done

# Modules that implement no genesis interfaces are silently skipped by the
# module manager's export (on this chain: consensus, MT, NFT, params — all
# with a null default). They fall back to their fresh default, but only when
# that default is an empty genesis (null / {}); a module with real expected
# data producing no export output still fails loudly.
for m in "${MODULES[@]}"; do
  if [[ -z ${FRAG_FILE[$m]:-} ]]; then
    d=$(tr -d ' \t\n\r' < "${DEF_FILE[$m]}")
    if [[ $d == 'null' || $d == '{}' ]]; then
      printf 'notice: module %s exports no genesis; using fresh default\n' "$m"
    else
      die "module $m produced no export output but has a non-empty default genesis"
    fi
  fi
done

ALL_MERGED=$(printf '%s\n' "${MODULES[@]}" "${RESET_MODULES[@]}" | sort -u)

### bank denom_metadata sanitize ############################################

# The legacy chain's bank denom_metadata violates several v0.50 validation
# rules: blank name/symbol (the token module never fills them), duplicate
# denom units, and IBC entries whose first denom unit is the source denom
# instead of the ibc/... base. name/symbol/display are display-only fields,
# so fill blanks from the token module's exported genesis (matched by
# min_unit == base), falling back to description/display; drop duplicate
# units; rebuild unusable unit lists as [base:0] and fall back display to
# base when the display cannot be expressed as a denom unit. Every change is
# logged.
if [[ -n ${FRAG_FILE[bank]:-} ]]; then
  BANK_FRAG=${FRAG_FILE[bank]}
  TOKEN_FRAG=${FRAG_FILE[token]:-}
  json_tool top "$BANK_FRAG" denom_metadata > "$WORK/denom-metadata.json" \
    || die "failed to extract denom_metadata from the bank fragment"
  python3 - "$WORK/denom-metadata.json" "$TOKEN_FRAG" \
    "$WORK/denom-metadata-fixed.json" > "$WORK/sanitize-count.txt" <<'PY2'
import json, sys

arr = json.load(open(sys.argv[1]))
tokens = {}
if len(sys.argv) > 2 and sys.argv[2]:
    tg = json.load(open(sys.argv[2]))
    for t in tg.get('tokens', []):
        tokens[t.get('min_unit')] = t

changes = 0
for m in arr:
    t = tokens.get(m.get('base')) or {}
    if not str(m.get('name', '')).strip():
        m['name'] = (str(t.get('name', '')).strip()
                     or str(m.get('description', '')).strip()
                     or m.get('base'))
        changes += 1
        print('denom_metadata: %s: name <- %r' % (m['base'], m['name']), file=sys.stderr)
    if not str(m.get('symbol', '')).strip():
        m['symbol'] = (str(t.get('symbol', '')).strip()
                       or m.get('display')
                       or m.get('base'))
        changes += 1
        print('denom_metadata: %s: symbol <- %r' % (m['base'], m['symbol']), file=sys.stderr)
    # v0.50 requires the first denom unit to be the base denom (exponent 0);
    # legacy IBC entries only carry a single exponent-0 unit named after the
    # source denom, which cannot be expressed under the new rules — rebuild
    if not m.get('denom_units') or m['denom_units'][0].get('denom') != m.get('base'):
        old = [(u.get('denom'), u.get('exponent')) for u in m.get('denom_units', [])]
        m['denom_units'] = [{'denom': m['base'], 'exponent': 0, 'aliases': []}]
        changes += 1
        print('denom_metadata: %s: rebuilt denom_units to [base:0] (legacy units: %s)'
              % (m['base'], old), file=sys.stderr)
    seen = set()
    units = []
    for u in m.get('denom_units', []):
        if u.get('denom') in seen:
            changes += 1
            print('denom_metadata: %s: dropped duplicate denom unit %r (exponent %s)'
                  % (m['base'], u.get('denom'), u.get('exponent')), file=sys.stderr)
            continue
        seen.add(u.get('denom'))
        units.append(u)
    m['denom_units'] = units
    # display must be one of the denom units; trace-path displays cannot be
    # units (they would need a fabricated exponent), so fall back to base
    if not any(u.get('denom') == m.get('display') for u in m.get('denom_units', [])):
        changes += 1
        print('denom_metadata: %s: display %r -> base %r (not expressible as a denom unit)'
              % (m['base'], m.get('display'), m['base']), file=sys.stderr)
        m['display'] = m['base']

json.dump(arr, open(sys.argv[3], 'w'))
print(changes)
PY2
  SANITIZED=$(cat "$WORK/sanitize-count.txt")
  if [[ $SANITIZED != 0 ]]; then
    json_tool patch "$BANK_FRAG" "$BANK_FRAG.new" denom_metadata \
      "$WORK/denom-metadata-fixed.json" || die "failed to patch bank fragment"
    mv "$BANK_FRAG.new" "$BANK_FRAG"
    printf 'bank denom_metadata sanitized: %s change(s)\n' "$SANITIZED"
  fi
fi

### gov params sanitize ######################################################

# Legacy gov params can carry expedited_min_deposit <= min_deposit, which
# the v0.50 genesis validation rejects (it must be strictly greater). Reset
# it to 5x min_deposit — the SDK's own default ratio (fresh init uses
# 10000000/50000000 uiris) — so expedited proposals cost 5x a regular
# deposit until governance changes the param.
if [[ -n ${FRAG_FILE[gov]:-} ]]; then
  GOV_FRAG=${FRAG_FILE[gov]}
  json_tool top "$GOV_FRAG" params > "$WORK/gov-params.json" \
    || die "failed to extract params from the gov fragment"
  python3 - "$WORK/gov-params.json" "$WORK/gov-params-fixed.json" \
    > "$WORK/gov-count.txt" <<'PY3'
import json, sys

p = json.load(open(sys.argv[1]))
fixed = 0
min_list = p.get('min_deposit') or []
if min_list:
    min_d = min_list[0]
    old = p.get('expedited_min_deposit') or []
    exp_d = old[0] if old else {}
    try:
        violates = (exp_d.get('denom') != min_d.get('denom')
                    or int(exp_d.get('amount', 0)) <= int(min_d.get('amount', 0)))
    except (TypeError, ValueError):
        violates = True
    if violates:
        new = {'denom': min_d['denom'], 'amount': str(int(min_d['amount']) * 5)}
        p['expedited_min_deposit'] = [new]
        fixed = 1
        print('gov params: expedited_min_deposit %s -> %s (must be strictly '
              'greater than min_deposit %s; using 5x like the SDK default)'
              % (json.dumps(old), json.dumps(new),
                 min_d['amount'] + min_d['denom']), file=sys.stderr)
json.dump(p, open(sys.argv[2], 'w'))
print(fixed)
PY3
  GOV_FIXED=$(cat "$WORK/gov-count.txt")
  if [[ $GOV_FIXED != 0 ]]; then
    json_tool patch "$GOV_FRAG" "$GOV_FRAG.new" params "$WORK/gov-params-fixed.json" \
      || die "failed to patch gov fragment"
    mv "$GOV_FRAG.new" "$GOV_FRAG"
    printf 'gov params sanitized: expedited_min_deposit reset to 5x min_deposit\n'
  fi
fi

### tibc client state sanitize ###############################################

# tibc-go v0.6.0 rejects legacy tibc genesis data on three counts: client
# states with trusting_period >= unbonding_period (both were created equal,
# e.g. 2400h), consensus states whose light client type differs from the
# client state type (relics from before a client type change), and consensus
# states / metadata belonging to clients no longer in the genesis. Fix the
# trusting periods to 2/3 of the unbonding period (the comet convention) and
# drop the invalid consensus states / metadata.
if [[ -n ${FRAG_FILE[tibc]:-} ]]; then
  TIBC_FRAG=${FRAG_FILE[tibc]}
  json_tool top "$TIBC_FRAG" client_genesis > "$WORK/tibc-clients.json" \
    || die "failed to extract client_genesis from the tibc fragment"
  python3 - "$WORK/tibc-clients.json" "$WORK/tibc-clients-fixed.json" \
    > "$WORK/tibc-count.txt" <<'PY4'
import json, re, sys


def parse_dur(s):
    m = re.fullmatch(r'(?:(\d+)h)?(?:(\d+)m)?(?:(\d+(?:\.\d+)?)s)?', s or '')
    if not m:
        return None
    return (int(m.group(1) or 0) * 3600 + int(m.group(2) or 0) * 60
            + float(m.group(3) or 0))


def fmt_dur(sec):
    sec = int(sec)
    return '%dh%dm%ds' % (sec // 3600, (sec % 3600) // 60, sec % 60)


g = json.load(open(sys.argv[1]))
fixed = 0
for c in g.get('clients') or []:
    cs = c.get('client_state') or {}
    tp = parse_dur(cs.get('trusting_period', ''))
    up = parse_dur(cs.get('unbonding_period', ''))
    if tp is not None and up is not None and tp >= up:
        new = fmt_dur(up * 2 // 3)
        print('tibc client %s: trusting_period %s -> %s (must be < unbonding_period %s)'
              % (c.get('chain_name'), cs['trusting_period'], new,
                 cs['unbonding_period']), file=sys.stderr)
        cs['trusting_period'] = new
        fixed += 1

# consensus states must be of the same light client type as the client state
# (legacy clients can carry states from before a client type change) and must
# belong to a client in the genesis; drop the invalid ones
def lc_package(type_url):
    m = re.search(r'lightclients\.([a-z0-9-]+)\.', type_url or '')
    return m.group(1) if m else None

client_pkg = {}
for c in g.get('clients') or []:
    client_pkg[c.get('chain_name')] = lc_package(
        (c.get('client_state') or {}).get('@type', ''))

kept_cc = []
for cc in g.get('clients_consensus') or []:
    name = cc.get('chain_name')
    if name not in client_pkg:
        print('tibc: dropped consensus states of orphan client %s' % name, file=sys.stderr)
        fixed += 1
        continue
    cpkg = client_pkg[name]
    kept_states = []
    for s in cc.get('consensus_states') or []:
        spkg = lc_package((s.get('consensus_state') or {}).get('@type', ''))
        if spkg and cpkg and spkg != cpkg:
            h = (s.get('height') or {}).get('revision_height')
            print('tibc client %s: dropped consensus state at height %s (type %s does not match client state type %s)'
                  % (name, h, spkg, cpkg), file=sys.stderr)
            fixed += 1
            continue
        kept_states.append(s)
    cc['consensus_states'] = kept_states
    kept_cc.append(cc)
g['clients_consensus'] = kept_cc

kept_md = []
for md in g.get('clients_metadata') or []:
    if md.get('chain_name') not in client_pkg:
        print('tibc: dropped metadata of orphan client %s' % md.get('chain_name'), file=sys.stderr)
        fixed += 1
        continue
    kept_md.append(md)
g['clients_metadata'] = kept_md

print('tibc sanitize: %d client(s), %d consensus entr(y/ies), %d consensus state(s) checked, %d fix(es)'
      % (len(g.get('clients') or []),
         len(g.get('clients_consensus') or []),
         sum(len(cc.get('consensus_states') or [])
             for cc in g.get('clients_consensus') or []),
         fixed), file=sys.stderr)

json.dump(g, open(sys.argv[2], 'w'))
print(fixed)
PY4
  TIBC_FIXED=$(cat "$WORK/tibc-count.txt")
  if [[ $TIBC_FIXED != 0 ]]; then
    json_tool patch "$TIBC_FRAG" "$TIBC_FRAG.new" client_genesis \
      "$WORK/tibc-clients-fixed.json" || die "failed to patch tibc fragment"
    mv "$TIBC_FRAG.new" "$TIBC_FRAG"
    printf 'tibc client states sanitized: %s client(s) fixed\n' "$TIBC_FIXED"
  fi
fi

### padded validation #######################################################

# validate every module and report ALL failures at once (a per-module die
# would hide failures behind the first one)
if [[ $SKIP_VALIDATE != 1 ]]; then
  FAILED_MODULES=()
  for m in "${EXPORTED[@]}"; do
    {
      printf '{'
      first=1
      for d in $ALL_MERGED; do
        if ((first)); then first=0; else printf ','; fi
        printf '"%s":' "$d"
        if [[ $d == "$m" ]]; then
          cat "${FRAG_FILE[$d]}"
        else
          cat "${DEF_FILE[$d]}"
        fi
      done
      printf '}'
    } > "$WORK/padded-app-state.json"
    json_tool skeleton "$FRESH" "$WORK/padded-genesis.json" \
      "$NEW_CHAIN_ID" "$WORK/padded-app-state.json" \
      || die "padded genesis assembly failed for module $m"
    if "$IRIS_BIN" genesis validate "$WORK/padded-genesis.json" \
      > "$WORK/validate-$m.log" 2>&1; then
      printf 'validated module: %s\n' "$m"
    else
      printf 'validation FAILED: module %s\n' "$m"
      grep -m1 '^Error:' "$WORK/validate-$m.log" | sed 's/^/  /'
      FAILED_MODULES+=("$m")
    fi
    rm -f "$WORK/padded-genesis.json" "$WORK/padded-app-state.json"
  done
  if ((${#FAILED_MODULES[@]} > 0)); then
    die "padded validation failed for module(s): ${FAILED_MODULES[*]} (fix the data and re-run with RESUME=1)"
  fi
fi

### final assembly ##########################################################

{
  printf '{'
  first=1
  for m in $ALL_MERGED; do
    if ((first)); then first=0; else printf ','; fi
    printf '"%s":' "$m"
    if [[ -n ${FRAG_FILE[$m]:-} ]]; then
      cat "${FRAG_FILE[$m]}"
    else
      cat "${DEF_FILE[$m]}"
    fi
  done
  printf '}'
} > "$WORK/app_state.json"

first_partial=$(ls "$WORK"/partial-*.json | head -n 1)
json_tool skeleton "$first_partial" "$OUT_DIR/genesis.json" \
  "$NEW_CHAIN_ID" "$WORK/app_state.json" || die "genesis assembly failed"

# the legacy chain carries evm code/storage entries for addresses that have
# no auth account (unreachable orphan state); evm InitGenesis panics on the
# missing account — drop them from the assembled genesis
json_tool filter-evm "$OUT_DIR/genesis.json" "$OUT_DIR/genesis.json.tmp" \
  || die "evm accounts sanitize failed"
mv "$OUT_DIR/genesis.json.tmp" "$OUT_DIR/genesis.json"

# legacy chains can carry a drift between the staking module's cached
# total_liquid_staked_tokens counter and the sum of the validators'
# liquid_shares; the crisis module's liquid stake invariant panics on it at
# genesis import — normalize the counter to the sum
json_tool fix-liquid-stake "$OUT_DIR/genesis.json" "$OUT_DIR/genesis.json.tmp" \
  || die "liquid staking counter sanitize failed"
mv "$OUT_DIR/genesis.json.tmp" "$OUT_DIR/genesis.json"

### assertions ###############################################################

cid=$(json_tool top "$OUT_DIR/genesis.json" chain_id) \
  || die "failed to read chain_id from $OUT_DIR/genesis.json"
[[ $cid == "\"$NEW_CHAIN_ID\"" ]] \
  || die "genesis chain_id=$cid, want \"$NEW_CHAIN_ID\""

ih=$(json_tool top "$OUT_DIR/genesis.json" initial_height) \
  || die "failed to read initial_height from $OUT_DIR/genesis.json"
[[ $ih == "$EXPECTED_INITIAL_HEIGHT" ]] \
  || die "genesis initial_height=$ih, want $EXPECTED_INITIAL_HEIGHT"

EXPORTED_SUPPLY=$(json_tool supply "$OUT_DIR/genesis.json" uiris) \
  || die "uiris not found in genesis bank supply"
[[ $EXPORTED_SUPPLY == "$EXPECTED_SUPPLY" ]] \
  || die "uiris supply=$EXPORTED_SUPPLY, want $EXPECTED_SUPPLY"

printf 'genesis: %s\nuiris supply: %s\n' \
  "$OUT_DIR/genesis.json" "$EXPORTED_SUPPLY"
