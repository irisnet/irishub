## IRIShub 简介

IRIS Hub是在Cosmos生态中的区域性枢纽，提供iService服务

## 现有资源

* 源码：https://github.com/irisnet/irishub
* 浏览器：
* 水龙头：
* 文档：https://github.com/irisnet/testnets



## 安装IRIShub

### 方法1：下载发行版安装

进入下载页: https://github.com/irisnet/irishub/releases/
下载对应版本的可执行文件
解压缩tar -C /usr/local/bin -xzf 、irishub_0.2.0_linux_amd64.zip
拷贝到/usr/local/bin/目录下 
执行以下命令,若出现对应的版本号则说明安装成功。
```
$ iris version
v0.2.0
    
$ iriscli version
v0.2.0
```
### 方法2：源码编译安装

### 安装Go版本 1.10+ 


系统要求：

Ubuntu LTS 16.04


安装IRISHub需要保证Go的版本在1.10以上，

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
    echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
    echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
    source ~/.bash_profile
```

参考链接：

1. https://golang.org/doc/install
2. https://github.com/golang/go/wiki/Ubuntu



### 下载源码并安装


在完成Go的安装后，通过以下命令下载并安装IRIS hub相关程序.

```
    mkdir -p $GOPATH/src/github.com/irisnet
    cd $GOPATH/src/github.com/irisnet
    git clone https://github.com/irisnet/irishub
    cd irishub && git checkout v0.2.0

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

    make get_vendor_deps && make install
```

以上命令将完成 iris 和 iriscli的安装. 若出现对应的版本号则说明安装成功。

```
    $ iris version
    v0.2.0
    
    $ iriscli version
    v0.2.0
```

