# IRIS Hub 网络迁移指南（irishub-1 → irishub-2）

本文档面向验证人，说明如何从 irishub-1 的停链快照迁移到新链 irishub-2。
新网络基于高度 37242247 的状态导出构建，**验证人共识私钥无需更换**。

**新网络统一使用二进制版本 `v4.1.0`**（官方发布的迁移版构建，包含导出补丁
与新链所需的全部改动；全体验证人请使用同一版本）。

官方 genesis 下载地址：

```
https://irishub-snapshot.oss-ap-southeast-1.aliyuncs.com/genesis.json
```

文件校验值（约 2.2GB；请同时以官方公告渠道公布的为准）：

```
MD5:    199c0f48900cc00a71c68e7c44a80c5d
SHA256: 6540bd7c1d11d5475e7eaabe8d16278c76dc319636cb1679a1f1351646822ef2
```

---

## 1. 快速迁移：下载官方 genesis 并启动（推荐）

### 1.1 前置条件

| 项目 | 要求 |
|---|---|
| 二进制 | **irishub `v4.1.0`**（官方发布的迁移版构建，含导出补丁） |
| 磁盘 | 本地 SSD，≥ 50GB 空闲（genesis 2.2GB + 初始状态 ~4-6GB + 增长余量） |
| 内存 | 启动日建议 ≥ 32GB（详见第 3 节）；16GB 需按第 3 节调参 |
| 验证人身份 | irishub-1 的 `priv_validator_key.json`（共识私钥，**不换**） |

### 1.2 下载并校验

```bash
wget -c https://irishub-snapshot.oss-ap-southeast-1.aliyuncs.com/genesis.json
shasum -a 256 genesis.json
# 必须等于: 6540bd7c1d11d5475e7eaabe8d16278c76dc319636cb1679a1f1351646822ef2
```

> `wget -c` 支持断点续传。大文件传输务必核对校验值后再使用。

### 1.3 准备节点目录

```bash
# 新 home（不要复用旧链目录）
iris init <你的moniker> --chain-id irishub-2 --home /data/irishub-2

# 覆盖 genesis
cp genesis.json /data/irishub-2/config/genesis.json

# 恢复验证人身份（关键步骤：共识私钥决定你在新链上的验证人席位）
cp <旧链home>/config/priv_validator_key.json /data/irishub-2/config/
cp <旧链home>/config/node_key.json           /data/irishub-2/config/   # 保留 P2P 身份（可选）

# 确认 genesis 与身份就位
iris genesis validate /data/irishub-2/config/genesis.json
```

**说明**：新链 genesis 的验证人集合（38 个 bonded 验证人）与 irishub-1 完全一致，
你的 `priv_validator_key.json` 对应的共识地址在其中，席位、投票权、委托关系全部原样继承。
启动前**不要**修改 staking 相关的任何字段。

### 1.4 启动（含内存调参）

```bash
GOMEMLIMIT=12GiB GOGC=30 iris start --home /data/irishub-2
```

16GB 机器还需（一次性，详见第 3 节）：

```bash
cat /proc/sys/vm/swappiness          # 若为 0，执行下行
sysctl -w vm.swappiness=60
systemctl stop systemd-oomd 2>/dev/null && systemctl mask systemd-oomd   # Ubuntu 22.04+
```

### 1.5 成功标志

```
INF InitChain chainID=irishub-2 initialHeight=37242248 module=server
INF initializing blockchain state from genesis.json module=server
（数分钟到数十分钟，取决于内存/磁盘，期间机器卡顿属正常，不要中断）
INF committed state ... height=37242248
```

之后节点正常出块（高度从 37242248 起）。可用以下命令核对：

```bash
iris q bank total              # 首块后应≈ 2156713998266827uiris（此后每块按 mint 参数增发）
iris q staking validators       # 应看到 38 个 bonded 验证人
```

---

## 2. 自行生成 genesis.json（可选，用于独立校验）

若希望不依赖官方分发的文件、用自己保存的链数据重新导出，使用仓库中的脚本：

```
scripts/chain-recovery/export-genesis.sh
```

### 2.1 前置条件

- **irishub-1 数据副本**：包含高度 37242247 状态的节点数据目录（即停链前的完整
  数据，2TB 量级；归档节点或停链快照均可）；
- **迁移版二进制 `v4.1.0`**：必须是包含导出补丁的构建（`appExport` 中禁用
  IAVL fastnode 自动迁移，否则 2TB 数据会触发全量 fastnode 重写）。官方发布的
  `v4.1.0` 已包含；自行编译：`git checkout v4.1.0 && go build -o iris ./cmd/iris`；
- 依赖：`bash`、`jq`、`python3`。

### 2.2 运行

```bash
IRIS_BIN=/path/to/iris \
OLD_HOME=/path/to/stopped-node-copy \
OUT_DIR=/path/to/export-output \
bash scripts/chain-recovery/export-genesis.sh
```

关键环境变量（默认值即本次迁移使用的官方值）：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `EXPORT_HEIGHT` | `37242247` | 导出所依据的链高度 |
| `NEW_CHAIN_ID` | `irishub-2` | 新链 ID |
| `EXPECTED_INITIAL_HEIGHT` | `37242248` | 逐个 partial 断言初始高度 |
| `EXPECTED_SUPPLY` | `2156713998266827` | 最终 genesis 的 uiris 供应量断言 |
| `RESUME` | `0` | 中断续跑置 1（按内容匹配已有 partial，不重导） |
| `MODULE_GROUPS` | 每模块一组 | 换行分隔的逗号连接模块组，可合并小模块减少轮次 |
| `SKIP_VALIDATE` | `0` | 内存极度紧张时置 1 跳过逐模块校验（断言仍执行） |

脚本流程：逐模块导出（内存上限 = 最大单模块，避免整链 OOM）→ 历史数据
sanitize（见 4.3）→ 逐模块校验 → 流式拼装最终 genesis（O(1) 内存）→
断言 `initial_height` 与 `uiris` 供应量。

成功标志（最后一行输出）：

```
genesis: /path/to/export-output/genesis.json
uiris supply: 2156713998266827
```

### 2.3 注意

- 全程需要对 2TB 数据做约 33 轮读取，耗时数小时；中断后 `RESUME=1` 续跑；
- 自行生成的 genesis 与官方文件**语义等价**（脚本内置供应量/高度断言保证经济
  数据一致），但字节级不完全相同（sanitize 日志量、空白符等），MD5 不同属预期；
- 可用 `scripts/chain-recovery/simulate-launch.sh` 在本地用自有验证人替换
  genesis 验证人集合，做启动前功能演练（出块、转账验证）。

---

## 3. 启动阶段内存需求与故障排查

### 3.1 内存峰值来自哪里

首次启动分两段，实测（2.2GB genesis）：

| 阶段 | 18GB 机器（调参后） | 32GB 机器（预估） |
|---|---|---|
| InitChain 状态导入 | ~5 分钟 | 更快 |
| **首个 Commit（状态树刷盘）** | ~20 分钟（内存吃满，靠 swap/压缩磨） | ~10 分钟内 |
| 之后出块 | 正常速度 | 正常速度 |

峰值构成：genesis 文件字节（2.2GB，全程持有）+ 状态树物化（~10-12GB）+
GC 与写缓冲余量 ≈ **15-18GB**。这是**一次性成本**：Commit 完成后状态落盘，
日常运行回到常规开销（irishub-1 验证人以小配置运行了七年）。

### 3.2 建议

- **推荐**：启动日使用 ≥32GB 内存的机器（可按量计费租用数小时，费用几十元
  级），完成首 Commit、出块稳定后，将 `data/` 目录迁回原机器继续运行；
- 16GB 机器务必按 1.4 调参（`GOMEMLIMIT=12GiB GOGC=30` + `swappiness=60` +
  关闭 systemd-oomd），可完成但耗时长（可能 1 小时+），期间机器近乎无响应；
- 必须使用本地 SSD，避免网络盘（首 Commit 要在内存压力下写 ~3GB 状态）。

### 3.3 故障排查

| 现象 | 原因 | 处置 |
|---|---|---|
| 进程被 `Killed`，swap 几乎没用 | 内核 OOM：常见 `vm.swappiness=0` 导致内核拒绝换页，撞墙即杀 | `dmesg -T \| grep -i oom` 确认；`sysctl -w vm.swappiness=60` 后重试 |
| 进程被杀，日志有 systemd-oomd | Ubuntu 22.04+ 的 oomd 按内存压力杀进程 | `systemctl stop systemd-oomd && systemctl mask systemd-oomd`（启动完成后恢复） |
| `InitChain` 后长时间无区块 | 首个 Commit 在刷盘，**正常现象** | 观察 `du -s <home>/data` 持续增长即是在推进，耐心等待，**不要杀进程** |
| 启动失败需要重试 | 中断留下半成品状态 | `rm -rf <home>/data`（**只删 data，保留 config**）后重新 start |
| 校验和不对 | 下载损坏 | `wget -c` 断点续传重下，再核对 SHA256 |
| 机器卡死、SSH 卡顿 | 内存压力传导到全系统 | 属预期；等待 Commit 完成；勿在此时强制重启 |

---

## 4. 保留的核心数据（模块清单）

### 4.1 完整保留（含全部账户与资产数据）

| 模块 | 数据量 | 内容 |
|---|---|---|
| evm | 1540 MB | EVM 账户、合约代码与存储 |
| nft | 390 MB | NFT 全量数据 |
| auth | 88 MB | 全部账户 |
| bank | 52 MB | 余额、供应量、代币元数据 |
| distribution | 17 MB | 委托奖励记录 |
| staking | 6 MB | 验证人（38 bonded + 269 unbonded）、委托、解委托、再委托、**3 条 tokenize-share 记录** |
| feegrant / slashing / 其余各模块 | <5 MB | gov、token、service、htlc、record、oracle、random、farm、coinswap、mt、vesting、upgrade 等 |

### 4.2 重置为初始默认（跨链状态不迁移）

`07-tendermint`、`ibc`、`transfer`、`interchainaccounts`、
`nonfungibletokentransfer`、`tibc`

> 所有 IBC/TIBC 客户端、连接、通道在新链上从零开始，相关业务方需要重新建立
> 跨链连接。此数据未被丢弃，可在旧链归档节点查证（见第 5 节）。

### 4.3 历史数据修订（为满足新版本校验所做的最小化修正，均有日志留痕）

| 修订 | 内容 |
|---|---|
| bank denom_metadata | 填充空的 name/symbol、重建不合规则的 denom_units、display 回退（共 153 处，展示字段，不影响余额） |
| gov 参数 | `expedited_min_deposit` 由 50 iris 调整为 25000 iris（=5× min_deposit 5000 iris，新版要求前者严格大于后者） |
| evm 账户过滤 | 剔除无对应 auth 账户的 evm 条目（不可达的死数据） |
| staking 计数器 | `total_liquid_staked_tokens` 与验证人 liquid_shares 求和对齐（1 单位级历史漂移） |

`consensus`、`MT`、`NFT`、`params` 四个模块本身无 genesis 导出接口，以默认
空值进入新链（与旧链实际语义一致）。

---

## 5. 原网络（irishub-1）的保留安排

官方将在迁移后**保留 irishub-1 全节点（含归档 RPC）一段时间**（至少 6 个月，
具体下线时间以官方公告为准），用途：

- 查询旧链全部历史：交易、区块、事件、任意高度的状态；
- 为被重置的数据（IBC/TIBC 客户端与通道历史、被过滤的 evm 死数据）**保留
  历史证明**；
- 若有个别资产在新链出现争议或遗漏，可回旧链核对原始状态。

同时社区第三方归档 RPC（如 `https://mainnet-iris-rpc.konsortech.xyz/`，支持
任意高度查询）也可用于交叉验证。irishub-1 停链后为只读状态，不再出块。

---

## 6. 迁移后网络基本情况

| 项目 | 值 |
|---|---|
| 二进制版本 | `v4.1.0`（全体验证人统一使用；genesis 内 `app_version` 字段为导出时的构建号 `4.0.3-9-g0e5e8c417`，仅元数据、不影响运行） |
| chain-id | `irishub-2` |
| 起始区块高度 | **37242248**（首块即此高度，延续旧链高度） |
| genesis 供应量 | `2156713998266827 uiris`（与旧链断点精确一致；此后每块按 mint 参数增发） |
| 验证人集合 | 38 个 bonded 验证人原样继承（出块需 ≥2/3 投票权在线）；另有 269 个 unbonded 历史记录随数据保留 |
| 验证人私钥 | **不变**（沿用 irishub-1 的 `priv_validator_key.json`） |
| 跨链 | IBC / TIBC 状态清零，所有连接需重新建立 |
| 流动性质押 | 3 条 tokenize-share 记录原样保留 |
| gov 参数 | `expedited_min_deposit` 调整为 25000 iris（见 4.3） |
| 代币元数据 | 展示字段规范化（见 4.3），余额与供应量不受影响 |
