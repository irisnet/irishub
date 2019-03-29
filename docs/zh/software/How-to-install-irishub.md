# 如何安装IRIShub
有以下两种方式：
### 方法1：下载发行版安装

进入下载页: https://github.com/irisnet/irishub/releases/
下载对应版本的可执行文件
解压缩
```
unzip -C /usr/local/bin  iris$VERSION.$OS-$ARCH.zip
```
拷贝到/usr/local/bin/目录下 
执行以下命令,若出现对应的版本号则说明安装成功。

```
$ iris version
0.13.1-a4a738e-0
    
$ iriscli version
0.13.1-a4a738e-0
```
### 方法2：源码编译安装

#### 安装Go版本 1.10+ 


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
请下载最新版本的代码编译（例如：`git checkout v0.13.1`），参考：https://github.com/irisnet/irishub/releases/
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
0.13.1-a4a738e-0
    
$ iriscli version
0.13.1-a4a738e-0
```
