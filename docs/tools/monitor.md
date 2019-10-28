---
order: 1
---

# IRIS Hub Monitor

## Introduction

IRIShub can report and serve the Prometheus metrics, which can be consumed by Prometheus collector(s).

This functionality is disabled by default.

To enable the Prometheus metrics, set `prometheus = true` in your config file(config.toml). Metrics will be served under /metrics on 26660 port by default. This port can be changed in the config file (`prometheus_listen_addr = ":26660"`).

## Metrics

Application metrics, namespace: `iris`

| **Name**                          | **Type** | **Tags**                       | **Description**                                              |
| --------------------------------- | -------- | ------------------------------ | ------------------------------------------------------------ |
| module_stake_bonded_token         | Gauge    | validator_address              | Total bonded token by validator                              |
| module_stake_loosen_token         | Gauge    |                                | Total loose tokens                                           |
| module_stake_burned_token         | Gauge    |                                | Total burned token                                           |
| module_stake_slashed_token        | Counter  | validator_address              | Total slashed token by validator                             |
| module_stake_jailed               | Gauge    | validator_address              | Jailed status by validator, either 0 (not jailed) or 1 (jailed) |
| module_stake_power                | Gauge    | validator_address              | Voting power by validator                                    |
| module_upgrade_upgrade            | Gauge    |                                | Whether new software needs to be installed, either 0 (no) or 1 (yes) |
| module_upgrade_signal             | Gauge    | validator_address, version     | Whether validator have run the new version software, either 0 (no) or 1 (yes) |
| module_service_active_requests    | Gauge    |                                | Number of active requests                                    |
| module_gov_parameter              | Gauge    | parameter_key                  | Parameter of governance                                      |
| module_gov_proposal_status        | Gauge    | proposal_id                    | Status of proposal, 0:DepositPeriod 1:VotingPeriod 2:Pass 3:Reject 4:Other |
| module_gov_vote                   | Gauge    | proposal_id, validator_address | Validator vote result of a proposal, 0:Yes 1:No 2:NoWithVeto 3:Abstain |
| module_distribution_community_tax | Gauge    | height                         | Community tax accumulation                                   |
| v0_invariant_failure              | Counter  | error                          | Invariant failure stats                                      |

Consensus metrics, namespace: `tendermint`

| **Name**                             | **Type**  | **Tags**         | **Description**                                              |
| ------------------------------------ | --------- | ---------------- | ------------------------------------------------------------ |
| consensus_height                     | Gauge     |                  | Height of the chain                                          |
| consensus_failure                    | Counter   | height           | Consensus failure                                            |
| consensus_validators                 | Gauge     |                  | Number of validators                                         |
| consensus_validators_power           | Gauge     |                  | Total voting power of all validators                         |
| consensus_missing_validators         | Gauge     |                  | Number of validators who did not sign                        |
| consensus_missing_validators_power   | Gauge     |                  | Total voting power of the missing validators                 |
| consensus_byzantine_validators       | Gauge     |                  | Number of validators who tried to double sign                |
| consensus_byzantine_validators_power | Gauge     |                  | Total voting power of the byzantine validators               |
| consensus_block_interval_seconds     | Histogram |                  | Time between this and last block (Block.Header.Time) in seconds |
| consensus_rounds                     | Gauge     |                  | Number of rounds                                             |
| consensus_num_txs                    | Gauge     |                  | Number of transactions                                       |
| consensus_block_parts                | Counter   | peer_id          | Number of blockparts transmitted by peer                     |
| consensus_latest_block_height        | Gauge     |                  | /status sync_info number                                     |
| consensus_fast_syncing               | Gauge     |                  | Either 0 (not fast syncing) or 1 (syncing)                   |
| consensus_total_txs                  | Gauge     |                  | Total number of transactions committed                       |
| consensus_block_size_bytes           | Gauge     |                  | Block size in bytes                                          |
| p2p_peers                            | Gauge     |                  | Number of peers node's connected to                          |
| p2p_peer_receive_bytes_total         | Counter   | peer_id          | Number of bytes received from a given peer                   |
| p2p_peer_send_bytes_total            | Counter   | peer_id          | Number of bytes sent to a given peer                         |
| p2p_peer_pending_send_bytes          | Gauge     | peer_id          | Number of pending bytes to be sent to a given peer           |
| p2p_num_txs                          | Gauge     | peer_id          | Number of transactions submitted by each peer_id             |
| mempool_size                         | Gauge     |                  | Number of uncommitted transactions                           |
| mempool_tx_size_bytes                | Histogram |                  | Transaction sizes in bytes                                   |
| mempool_failed_txs                   | Counter   |                  | Number of failed transactions                                |
| mempool_recheck_times                | Counter   |                  | Number of transactions rechecked in the mempool              |
| state_block_processing_time          | Histogram |                  | Time between BeginBlock and EndBlock in ms                   |
| state_recheck_time                   | Histogram |                  | Time cost on recheck in ms                                   |
| state_app_hash_conflict              | Counter   | proposer, height | App hash conflict error                                      |

IRIShub metrics also contains tendermint metrics, Visit [tendermint metrics](https://github.com/irisnet/tendermint/blob/irisnet/master/docs/tendermint-core/metrics.md) for more information.

## Start Monitor

This is an example for getting started with the IRIShub Monitor by using docker.

### Edit Prometheus config file

You can download the example [prometheus.yml](https://raw.githubusercontent.com/prometheus/prometheus/master/documentation/examples/prometheus.yml)  to the `~/volumes/prometheus/` and add a job under the `scrape_configs` of the `prometheus.yml`:

```yaml
      - job_name: irishub
          static_configs:
          - targets: ['localhost:26660']
            labels:
              instance: myvalidator
```

:::tip
The value of targets is the `prometheus_listen_addr` in your `config.toml`
:::

### Start Prometheus

```bash
docker run -d --name=prometheus -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus
```

You should be able to browse to a status page about prometheus at <http://localhost:9090>

### Start Grafana

```bash
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

You can visit <http://localhost:3000/> to open grafana and create your own dashboard.

:::tip
The default username and password are both admin. We strongly recommend immediately changing your username & password after login.

A Grafana dashboard compatible with all the cosmos-sdk and tendermint based blockchains: [cosmos-dashboard](https://github.com/zhangyelong/cosmos-dashboard)
:::
