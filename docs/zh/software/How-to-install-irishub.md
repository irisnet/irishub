# 如何安装`iris`

## 最新版本

当前IRIShub最新版本为 [v0.15.5](https://github.com/irisnet/irishub/releases/latest)

::: tip
请将下文中的 <latest_iris_version> 替换为 v0.15.5
:::

## 服务器配置要求

首先，你需要配置一台服务器。验证人节点要求一直稳定运行，所以你需要在一台数据中心的服务器上部署IRIShub。任何像AWS、GCP、DigitalOcean中的云服务器都是适合的。

**推荐的服务器的配置：**

- CPU核数：2
- 内存容量：6GB
- 磁盘空间：256GB SSD
- 操作系统：Ubuntu 16.04 LTS +
- 带宽: 20Mbps
- 允许的入方向的链接：TCP端口 26656 和 26657

## 安装

### 安装 `go`

::: tip
IRIShub需要 **Go 1.12.5+**.
:::

参照[官方文档](https://golang.org/doc/install)安装 `go`。

记住要设置 `$GOPATH`，`$GOBIN` 和 `$PATH` 环境变量，示例:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN" >> ~/.bashrc
source ~/.bashrc
```

确认Go是否安装成功：
```bash
go version
```

### 安装 `iris`

在完成Go的安装后，通过以下命令下载并安装IRIShub相关程序.

请确保你的电脑可以访问`google.com`， `iris`很多库的依赖由google提供（如果无法访问，你也可以尝试添加一个代理: `export GOPROXY=https://goproxy.io`）

```bash
git clone --branch <latest_iris_version> https://github.com/irisnet/irishub
cd irishub
# source scripts/setTestEnv.sh # 若参与测试网请执行此命令
make get_tools install
```

以上命令将完成`iris`和`iriscli`的安装. 若出现对应的版本号则说明安装成功。

```bash
$ iris version
<latest_iris_version>
    
$ iriscli version
<latest_iriscli_version>
```
