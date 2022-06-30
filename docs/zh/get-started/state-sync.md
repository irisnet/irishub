---
order: 7
---

# State Sync
## 背景

如果您在不需要历史数据的情况下要快速启动节点并加入 IRIS Hub，可以考虑使用 `state_sync` 功能快速启动节点。需要注意的是在启动节点的时候要保证节点 data 目录是空的。

## 步骤

1. 参考 [加入主网](./mainnet.md) 进行主网节点初始化
2. 查看现存快照的块高，并选取最新块高

```bash
curl http://34.82.96.8:26658/
```

3. 修改配置文件 `config.toml`

```toml
[statesync]
enable = true #是否开启 stat_sync,设置为 true
rpc_servers = "34.82.96.8:26657,34.77.68.145:26657"  #链接的 rpc server 地址
trust_height = # 设置为最新快照的块高
trust_hash = "" #设置为最新快照块高对应的 hash，可通过浏览器 https://irishub.iobscan.io/#/block/<trust_height> 进行查看
trust_period = "168h0m0s"
discovery_time = "15s"
temp_dir = ""
```

4. 启动节点

```bash
iris start
```

## 其他
1. 如果在启动链的过程中有问题可以执行 `iris unsafe-reset-all` 重置节点，然后重复以上步骤。
2. 如果出现不能解决的问题，请通过 [Discord](https://discord.com/invite/bmhu9F9xbX) 来联系 IRISnet，我们将帮助您解决遇到的问题。
