# 如何安装`iris`

### 当前IRIShub最新版本为 ： v0.13.1
获取最新版本 https://github.com/irisnet/irishub/releases/latest
```
注意：使用 git checkout 命令时，请将 <latest_iris_version> 替换成 v0.13.1
```

### 源码编译安装

#### 服务器配置要求

首先，你需要配置一台服务器。验证人节点要求一直稳定运行，所以你需要在一台数据中心的服务器上部署IRIShub。任何像AWS、GCP、DigitalOcean中的云服务器都是适合的。

IRIShub是用Go语言编写的。它可以在任何能够编译并运行Go语言程序的平台上工作。然而，强烈建议在Linux服务器上运行验证人节点。
推荐的服务器的配置：

* CPU核数：2
* 内存容量：6GB
* 磁盘空间：256GB SSD
* 操作系统：Ubuntu 18.04 LTS/16.04 LTS
* 带宽: 20Mbps
* 允许的入方向的链接：TCP端口 26656 和 26657


#### 安装Go版本 1.10+ 


系统要求：

Ubuntu LTS 16.04


安装IRIShub需要保证Go的版本在1.10以上，

通过执行以下命令安装1.10版本的Go。

```
$ sudo add-apt-repository ppa:gophers/archive
$ sudo apt-get update
$ sudo apt-get install golang-1.10-go
```

以上命令将安装 golang-1.10-go在 /usr/lib/go-1.10/bin. 需要将它加入到PATH中

```
echo "export PATH=$PATH:/usr/lib/go-1.10/bin" >> ~/.bash_profile
source ~/.bash_profile
```

同时，你需要指定相关的 $GOPATH, $GOBIN, 和 $PATH 变量, 例如:

```
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
source ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
source ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
source ~/.bash_profile
```

参考链接：

1. https://golang.org/doc/install
2. https://github.com/golang/go/wiki/Ubuntu



#### 下载源码并安装


在完成Go的安装后，通过以下命令下载并安装IRIS hub相关程序.(请确保你的电脑可以访问`google.com`， `iris`很多库的依赖由google提供)

* 编译用于`测试网`的可执行文件:
请下载最新版本的代码编译，参考：https://github.com/irisnet/irishub/releases/latest
```
mkdir -p $GOPATH/src/github.com/irisnet
cd $GOPATH/src/github.com/irisnet
git clone https://github.com/irisnet/irishub
cd irishub && git checkout <latest_iris_version>
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
make get_tools
make get_vendor_deps
source scripts/setTestEnv.sh
make all
```

* 编译用于`betanet`的可执行文件:
```
mkdir -p $GOPATH/src/github.com/irisnet
cd $GOPATH/src/github.com/irisnet
git clone https://github.com/irisnet/irishub
cd irishub && git checkout <latest_iris_version>
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
make get_tools
make get_vendor_deps
make all
```

以上命令将完成`iris`和`iriscli`的安装. 若出现对应的版本号则说明安装成功。

```
$ iris version
<latest_iris_version>
    
$ iriscli version
<latest_iriscli_version>
```
