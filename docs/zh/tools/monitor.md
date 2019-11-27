---
order: 1
---

# 监控工具

## 简介

IRIShub 软件可以提供Prometheus监控指标，Prometheus可以收集这些指标。

默认情况下此功能是被禁用的，要启用Prometheus监控指标，请在配置文件(config.toml)中设置`prometheus = true`。默认情况下，Metrics将在26660端口下的/metrics提供，可以在配置文件中修改服务地址`prometheus_listen_addr = ":26660"`。

## 监控指标

应用层监控参数如下，命名空间: `iris`

| Name                              | Type    | Tags                           | Description                                      |
| --------------------------------- | ------- | ------------------------------ | ------------------------------------------------ |
| module_stake_bonded_token         | Gauge   | validator_address              | 验证人被绑定的token总数                          |
| module_stake_loosen_token         | Gauge   |                                | 未被绑定的token总数                              |
| module_stake_burned_token         | Gauge   |                                | 销毁的token总数                                  |
| module_stake_slashed_token        | Counter | validator_address              | 验证人被惩罚的token总数                          |
| module_stake_jailed               | Gauge   | validator_address              | 验证人监禁状态，0（未监禁）或1（被监禁）         |
| module_stake_power                | Gauge   | validator_address              | 验证人投票权                                     |
| module_upgrade_upgrade            | Gauge   |                                | 是否需要安装新软件，0（否）或1（是）             |
| module_upgrade_signal             | Gauge   | validator_address, version     | 验证人是否运行了新版本软件，0（否）或1（是）     |
| module_service_active_requests    | Gauge   |                                | 活跃的请求数                                     |
| module_gov_parameter              | Gauge   | parameter_key                  | 治理参数                                         |
| module_gov_proposal_status        | Gauge   | proposal_id                    | 提议状态，0:抵押期 1:投票期 2:通过 3:拒绝 4:其他 |
| module_gov_vote                   | Gauge   | proposal_id, validator_address | 验证人投票结果，0:同意 1:反对 2:强烈反对 3:弃权  |
| module_distribution_community_tax | Gauge   | height                         | 社区基金累计值                                   |
| v0_invariant_failure              | Counter | error                          | Invariant检查错误事件                            |

共识层监控参数如下，名字空间：`tendermint`

| **Name**                             | **Type**  | **Tags**         | **Description**                                  |
| ------------------------------------ | --------- | ---------------- | ------------------------------------------------ |
| consensus_height                     | Gauge     |                  | 共识状态机所在高度                               |
| consensus_failure                    | Counter   | height           | 导致共识状态机终止的错误                         |
| consensus_validators_power           | Gauge     |                  | 验证人的投票权重之和                             |
| consensus_missing_validators         | Gauge     |                  | 块中缺失的precommit数量                          |
| consensus_missing_validators_power   | Gauge     |                  | 块中缺失的precommit对应的验证人所占的投票权重    |
| consensus_byzantine_validators       | Gauge     |                  | 块中包含的作恶证据对应的验证人数量               |
| consensus_byzantine_validators_power | Gauge     |                  | 作恶的验证人的投票权重                           |
| consensus_block_interval_seconds     | Histogram |                  | 从区块被提出到区块执行完毕消耗的时间             |
| consensus_rounds                     | Gauge     |                  | 共识状态机所处的round                            |
| consensus_num_txs                    | Gauge     |                  | 区块中包含的交易数量                             |
| consensus_num_txs                    | Gauge     |                  | 区块中包含的交易数量                             |
| consensus_block_parts                | Counter   | peer_id          | 区块被切分的块数                                 |
| consensus_latest_block_height        | Gauge     |                  | 共识状态机上一个高度                             |
| consensus_fast_syncing               | Gauge     |                  | 是否处于fast_sync模式                            |
| consensus_total_txs                  | Gauge     |                  | 区块链打包的交易总数                             |
| consensus_block_size_bytes           | Gauge     |                  | 区块大小                                         |
| p2p_peers                            | Gauge     |                  | 连接的peer的数量                                 |
| p2p_peer_receive_bytes_total         | Counter   | peer_id          | 从某个peer接受的字节数量                         |
| p2p_peer_send_bytes_total            | Counter   | peer_id          | 发给某个peer的字节数量                           |
| p2p_peer_pending_send_bytes          | Gauge     | peer_id          | 处于等待发送状态的字节数量                       |
| p2p_num_txs                          | Gauge     | peer_id          | 某个peer广播过来的交易数量                       |
| mempool_size                         | Gauge     |                  | mempool中交易数量                                |
| mempool_tx_size_bytes                | Histogram |                  | mempool中新增的交易大小                          |
| mempool_failed_txs                   | Counter   |                  | mempool收到的无法通过checkTx的交易数量           |
| mempool_recheck_times                | Counter   |                  | mempool中对多少交易执行过recheck                 |
| state_block_processing_time          | Histogram |                  | 交易执行所消耗的时间，不包含beginBlock和endBlock |
| state_recheck_time                   | Histogram |                  | Recheck消耗的时间                                |
| state_app_hash_conflict              | Counter   | proposer, height | AppHash冲突的错误                                |

IRIShub metrics也包含tendermint metrics，有关更多信息，请访问[tendermint metrics](https://github.com/irisnet/tendermint/blob/irisnet/master/docs/tendermint-core/metrics.md)。

## 启动监控工具

这是使用docker来启动IRIShub Monitor的示例。

### 编辑Prometheus配置文件

你可以将示例[prometheus.yml](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml)下载到`~/volumes/prometheus/`并在配置文件`prometheus.yml`中添加以下`jobs`：

```yaml
      - job_name: irishub
          static_configs:
          - targets: ['localhost:26660']
            labels:
              instance: myvalidator
```

:::tip
targets配置项的值为节点``config.toml`文件中`prometheus_listen_addr`的值。
:::

### 启动 Prometheus

```bash
docker run -d --name=prometheus -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus
```

你应该可以在<http://localhost:9090>上浏览到prometheus的状态页面。

### 启动 Grafana

```bash
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

你可以访问<http://localhost:3000/>打开Grafana并创建自己的仪表盘。

:::tip
默认的用户名和密码均为admin。强烈建议在登录后立即更改你的用户名和密码。

兼容所有基于 cosmos-sdk 和 tendermint 的区块链的 Grafana 仪表盘: [cosmos-dashboard](https://github.com/zhangyelong/cosmos-dashboard)
:::
