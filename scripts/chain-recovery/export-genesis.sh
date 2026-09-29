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
#      so peak RSS is bounded by the largest single module;
#   2. assembles the final genesis with a streaming JSON member scanner
#      (python, O(1) memory), producing the same result as the old jq merge:
#      chain_id replaced, IBC modules reset to fresh-init defaults, everything
#      else taken from the export at EXPORT_HEIGHT;
#   3. validates each exported module by building a "padded" genesis — the
#      fresh defaults with that one module replaced by the exported fragment
#      (`iris genesis validate` requires every module section to be present,
#      so a bare partial cannot be validated directly);
#   4. asserts chain_id / initial_height / uiris supply with streaming tools
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

IRIS_BIN=${IRIS_BIN:-/path/to/verified/iris}
OLD_HOME=${OLD_HOME:-/path/to/stopped-node-copy}
OUT_DIR=${OUT_DIR:-/path/to/export-output}

NEW_CHAIN_ID=${NEW_CHAIN_ID:-irishub-2}
EXPORT_HEIGHT=${EXPORT_HEIGHT:-37242247}
EXPECTED_INITIAL_HEIGHT=${EXPECTED_INITIAL_HEIGHT:-37242248}
EXPECTED_SUPPLY=${EXPECTED_SUPPLY:-2156713998266827}

IBC_MODULES=(07-tendermint ibc transfer interchainaccounts nonfungibletokentransfer)
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
        if mode == 'skeleton' and len(sys.argv) == 6:
            return cmd_skeleton(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5])
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

### module list (export everything except the IBC modules, which are reset
### to fresh-init defaults in the final genesis) ##############################

mapfile -t ALL_KEYS < <(jq -r '.app_state | keys[]' "$FRESH" | sort)
declare -A IBC_SET=()
for m in "${IBC_MODULES[@]}"; do IBC_SET[$m]=1; done

MODULES=()
for m in "${ALL_KEYS[@]}"; do
  if [[ -z ${IBC_SET[$m]:-} ]]; then
    MODULES+=("$m")
  fi
done
((${#MODULES[@]} > 0)) || die "no modules to export"

for m in "${MODULES[@]}" "${IBC_MODULES[@]}"; do
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

i=0
for g in "${EXPORT_GROUPS[@]}"; do
  [[ -n $g ]] || continue
  i=$((i+1))
  partial="$WORK/partial-$i.json"
  if [[ $RESUME == 1 && -s $partial ]]; then
    printf 'resume: reusing %s\n' "$partial"
  else
    printf 'export group %d: %s\n' "$i" "$g"
    "$IRIS_BIN" export \
      --home "$OLD_HOME" \
      --height "$EXPORT_HEIGHT" \
      --for-zero-height=false \
      --modules-to-export "$g" \
      --output-document "$partial"
  fi

  ih=$(json_tool top "$partial" initial_height) \
    || die "failed to read initial_height from $partial"
  [[ $ih == "$EXPECTED_INITIAL_HEIGHT" ]] \
    || die "$partial: initial_height=$ih, want $EXPECTED_INITIAL_HEIGHT"
done

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

if ! diff <(printf '%s\n' "${MODULES[@]}") <(printf '%s\n' "${EXPORTED[@]}"); then
  die "exported module set does not match the expected module list"
fi

# fresh defaults for every module: used to reset the IBC modules in the final
# genesis and to pad partials for `iris genesis validate`
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

ALL_MERGED=$(printf '%s\n' "${MODULES[@]}" "${IBC_MODULES[@]}" | sort -u)

### padded validation #######################################################

if [[ $SKIP_VALIDATE != 1 ]]; then
  for m in "${MODULES[@]}"; do
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
    "$IRIS_BIN" genesis validate "$WORK/padded-genesis.json" \
      || die "genesis validate failed for module $m"
    rm -f "$WORK/padded-genesis.json" "$WORK/padded-app-state.json"
    printf 'validated module: %s\n' "$m"
  done
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
