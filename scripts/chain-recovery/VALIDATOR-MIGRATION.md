# IRIS Hub Network Migration Guide (irishub-1 → irishub-2)

This document is for validators. It explains how to migrate from the halted
irishub-1 snapshot to the new chain, irishub-2.
The new network is built from a state export at height 37242247. **Validator
consensus private keys do not need to be replaced.**

**The new network uses binary version `v4.1.0` exclusively** (the official
migration build, which contains the export patch and all changes required for
the new chain; all validators must use the same version).

Official genesis download:

```
https://irishub-snapshot.oss-ap-southeast-1.aliyuncs.com/genesis.json
```

Checksums (approx. 2.2GB; always cross-check against the official announcement
channels as well):

```
MD5:    199c0f48900cc00a71c68e7c44a80c5d
SHA256: 6540bd7c1d11d5475e7eaabe8d16278c76dc319636cb1679a1f1351646822ef2
```

---

## 1. Fast migration: download the official genesis and start (recommended)

### 1.1 Prerequisites

| Item | Requirement |
|---|---|
| Binary | **irishub `v4.1.0`** (official migration build, includes the export patch) |
| Disk | Local SSD, ≥ 50GB free (genesis 2.2GB + initial state ~4-6GB + growth headroom) |
| Memory | ≥ 32GB recommended on launch day (see section 3); 16GB machines must apply the tuning in section 3 |
| Validator identity | irishub-1's `priv_validator_key.json` (consensus private key, **unchanged**) |

### 1.2 Download and verify

```bash
wget -c https://irishub-snapshot.oss-ap-southeast-1.aliyuncs.com/genesis.json
shasum -a 256 genesis.json
# must equal: 6540bd7c1d11d5475e7eaabe8d16278c76dc319636cb1679a1f1351646822ef2
```

> `wget -c` supports resuming interrupted downloads. For a file this large,
> always verify the checksum before use.

### 1.3 Prepare the node home

```bash
# new home (do not reuse the old chain directory)
iris init <your-moniker> --chain-id irishub-2 --home /data/irishub-2

# overwrite genesis
cp genesis.json /data/irishub-2/config/genesis.json

# restore validator identity (critical step: the consensus private key determines
# your validator seat on the new chain)
cp <old-chain-home>/config/priv_validator_key.json /data/irishub-2/config/
cp <old-chain-home>/config/node_key.json           /data/irishub-2/config/   # keep P2P identity (optional)

# confirm genesis and identity are in place
iris genesis validate /data/irishub-2/config/genesis.json
```

**Note**: the validator set in the new genesis (38 bonded validators) is
identical to irishub-1's. The consensus address matching your
`priv_validator_key.json` is included, and seats, voting power, and
delegations are all inherited as-is.
Do **not** modify any staking-related fields before launch.

### 1.4 Start (with memory tuning)

```bash
GOMEMLIMIT=12GiB GOGC=30 iris start --home /data/irishub-2
```

On 16GB machines, also apply the following (one-time, see section 3 for details):

```bash
cat /proc/sys/vm/swappiness          # if 0, run the next line
sysctl -w vm.swappiness=60
systemctl stop systemd-oomd 2>/dev/null && systemctl mask systemd-oomd   # Ubuntu 22.04+
```

### 1.5 Success criteria

```
INF InitChain chainID=irishub-2 initialHeight=37242248 module=server
INF initializing blockchain state from genesis.json module=server
(minutes to tens of minutes depending on memory/disk; the machine may appear
frozen during this time, which is normal — do not interrupt)
INF committed state ... height=37242248
```

Afterwards the node produces blocks normally (starting from height 37242248).
Verify with:

```bash
iris q bank total              # should be ≈ 2156713998266827uiris right after the first block (then grows per-block per mint params)
iris q staking validators       # should show 38 bonded validators
```

---

## 2. Generate genesis.json yourself (optional, for independent verification)

If you prefer not to rely on the officially distributed file and want to
re-export from your own copy of chain data, use the repository script:

```
scripts/chain-recovery/export-genesis.sh
```

### 2.1 Prerequisites

- **irishub-1 data copy**: a node data directory containing state at height
  37242247 (i.e. complete pre-halt data, ~2TB scale; an archive node or a halt
  snapshot both work);
- **migration binary `v4.1.0`**: must be a build containing the export patch
  (disables the automatic IAVL fastnode migration in `appExport`, otherwise
  2TB of data triggers a full fastnode rewrite). The official `v4.1.0`
  release already contains it; to build it yourself:
  `git checkout v4.1.0 && go build -o iris ./cmd/iris`;
- dependencies: `bash`, `jq`, `python3`.

### 2.2 Run

```bash
IRIS_BIN=/path/to/iris \
OLD_HOME=/path/to/stopped-node-copy \
OUT_DIR=/path/to/export-output \
bash scripts/chain-recovery/export-genesis.sh
```

Key environment variables (defaults are the official values used for this
migration):

| Variable | Default | Description |
|---|---|---|
| `EXPORT_HEIGHT` | `37242247` | chain height the export is based on |
| `NEW_CHAIN_ID` | `irishub-2` | new chain ID |
| `EXPECTED_INITIAL_HEIGHT` | `37242248` | asserted initial height of every partial |
| `EXPECTED_SUPPLY` | `2156713998266827` | asserted uiris supply of the final genesis |
| `RESUME` | `0` | set to 1 to resume an interrupted run (reuses existing partials matched by content, no re-export) |
| `MODULE_GROUPS` | one module per group | newline-separated, comma-joined module groups; merging small modules reduces rounds |
| `SKIP_VALIDATE` | `0` | set to 1 to skip per-module validation under extreme memory pressure (assertions still run) |

Script flow: per-module export (peak memory bounded by the largest single
module, avoiding whole-chain OOM) → legacy data sanitize (see 4.3) →
per-module validation → streaming assembly of the final genesis (O(1) memory)
→ assertions on `initial_height` and `uiris` supply.

Success criteria (last line of output):

```
genesis: /path/to/export-output/genesis.json
uiris supply: 2156713998266827
```

### 2.3 Notes

- The full run reads ~2TB of data in about 33 passes and takes several hours;
  resume interrupted runs with `RESUME=1`;
- A self-generated genesis is **semantically equivalent** to the official file
  (the built-in supply/height assertions guarantee identical economic data),
  but not byte-identical (sanitize log volume, whitespace, etc.), so a
  different MD5 is expected;
- `scripts/chain-recovery/simulate-launch.sh` can replace the genesis
  validator set with your own validators locally as a pre-launch functional
  rehearsal (block production, transfer verification).

---

## 3. Launch-day memory requirements and troubleshooting

### 3.1 Where the memory peak comes from

First launch has two phases, measured with a 2.2GB genesis:

| Phase | 18GB machine (tuned) | 32GB machine (estimated) |
|---|---|---|
| InitChain state import | ~5 minutes | faster |
| **First Commit (state tree flush)** | ~20 minutes (memory saturated, grinding through swap/compression) | under ~10 minutes |
| Block production afterwards | normal speed | normal speed |

Peak composition: genesis file bytes (2.2GB, held throughout) + materialized
state tree (~10-12GB) + GC and write-buffer headroom ≈ **15-18GB**. This is a
**one-time cost**: once the Commit lands on disk, routine operation returns to
normal overhead (irishub-1 validators ran on small machines for seven years).

### 3.2 Recommendations

- **Recommended**: use a ≥32GB machine on launch day (an hourly-billed rental
  for a few hours costs only tens of yuan); once the first Commit lands and
  block production is stable, move the `data/` directory back to the original
  machine and keep running there;
- 16GB machines must apply the section 1.4 tuning (`GOMEMLIMIT=12GiB GOGC=30`
  + `swappiness=60` + disable systemd-oomd). It works but takes longer
  (possibly 1 hour+), during which the machine is nearly unresponsive;
- Use a local SSD, never network storage (the first Commit writes ~3GB of
  state under memory pressure).

### 3.3 Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| Process `Killed`, swap barely used | Kernel OOM: commonly `vm.swappiness=0` makes the kernel refuse paging and kill on impact | confirm with `dmesg -T \| grep -i oom`; retry after `sysctl -w vm.swappiness=60` |
| Process killed, systemd-oomd in logs | Ubuntu 22.04+ oomd kills processes under memory pressure | `systemctl stop systemd-oomd && systemctl mask systemd-oomd` (restore after launch) |
| No blocks long after `InitChain` | First Commit is flushing to disk, **normal** | growing `du -s <home>/data` means progress; wait patiently, **do not kill the process** |
| Retry after failed launch | Interrupted run left half-written state | `rm -rf <home>/data` (**data only, keep config**) and start again |
| Checksum mismatch | Corrupt download | resume with `wget -c`, then re-verify SHA256 |
| Machine frozen, SSH laggy | Memory pressure spreading system-wide | expected; wait for the Commit to finish; do not force reboot |

---

## 4. Preserved core data (module inventory)

### 4.1 Fully preserved (all accounts and asset data)

| Module | Size | Contents |
|---|---|---|
| evm | 1540 MB | EVM accounts, contract code and storage |
| nft | 390 MB | Full NFT data |
| auth | 88 MB | All accounts |
| bank | 52 MB | Balances, supply, token metadata |
| distribution | 17 MB | Delegation reward records |
| staking | 6 MB | Validators (38 bonded + 269 unbonded), delegations, undelegations, redelegations, **3 tokenize-share records** |
| feegrant / slashing / all other modules | <5 MB | gov, token, service, htlc, record, oracle, random, farm, coinswap, mt, vesting, upgrade, etc. |

### 4.2 Reset to fresh defaults (cross-chain state is not migrated)

`07-tendermint`, `ibc`, `transfer`, `interchainaccounts`,
`nonfungibletokentransfer`, `tibc`

> All IBC/TIBC clients, connections, and channels start from zero on the new
> chain; integrators need to re-establish cross-chain connections. This data
> was not discarded — it remains verifiable on the old-chain archive nodes
> (see section 5).

### 4.3 Legacy data fixes (minimal changes required by new-version validation, all logged)

| Fix | Contents |
|---|---|
| bank denom_metadata | filled blank name/symbol, rebuilt non-conforming denom_units, display fallback (153 places total; display fields only, balances untouched) |
| gov params | `expedited_min_deposit` adjusted from 50 iris to 25000 iris (= 5× min_deposit of 5000 iris; the new version requires the former to be strictly greater) |
| evm account filter | dropped evm entries with no matching auth account (unreachable dead data) |
| staking counter | aligned `total_liquid_staked_tokens` with the sum of validator liquid_shares (1-unit-level historical drift) |

The `consensus`, `MT`, `NFT`, and `params` modules expose no genesis export
interface, so they enter the new chain with default empty values. Note the
case sensitivity: uppercase `MT` / `NFT` are the TIBC cross-chain transfer apps
(cleared together with `tibc` under section 4.2), **not** the native `mt` /
`nft` modules — native MT/NFT collections are fully preserved (see 4.1).
`consensus` carries no genesis state by design.

---

## 5. Retention of the original network (irishub-1)

After the migration, irishub-1 full nodes (including archive RPC) will be
kept for **a period of time** (at least 6 months; the exact sunset date will be
announced separately), serving to:

- query all old-chain history: transactions, blocks, events, and state at any
  height;
- preserve **historical proofs** for the reset data (IBC/TIBC client and
  channel history, filtered evm dead data);
- cross-check the original state on the old chain if any asset on the new
  chain is disputed or found missing.

Community third-party archive RPCs (e.g.
`https://mainnet-iris-rpc.konsortech.xyz/`, supporting queries at arbitrary
heights) can also be used for cross-verification. irishub-1 is read-only after
the halt and no longer produces blocks.

---

## 6. Post-migration network facts

| Item | Value |
|---|---|
| Binary version | `v4.1.0` (uniform for all validators; the `app_version` field inside genesis is the export-time build tag `4.0.3-9-g0e5e8c417`, metadata only, does not affect operation) |
| chain-id | `irishub-2` |
| First block height | **37242248** (first block is at this height, continuing old-chain height) |
| Genesis supply | `2156713998266827 uiris` (exactly matching the old chain at the breakpoint; grows per-block per mint params afterwards) |
| Validator set | 38 bonded validators inherited as-is (≥2/3 voting power must be online to produce blocks); 269 unbonded historical records preserved with the data |
| Validator keys | **unchanged** (reuse irishub-1's `priv_validator_key.json`) |
| Cross-chain | IBC / TIBC state wiped; all connections must be re-established |
| Liquid staking | 3 tokenize-share records preserved as-is |
| Gov params | `expedited_min_deposit` adjusted to 25000 iris (see 4.3) |
| Token metadata | display fields normalized (see 4.3); balances and supply unaffected |
