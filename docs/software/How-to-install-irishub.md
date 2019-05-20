# How to install `iris` 

### The Latest version of IRIShub : v0.13.1
refer to : https://github.com/irisnet/irishub/releases/latest
```
Please replace <latest_iris_version> with v0.13.1 while using "git checkout" 
```

You can download the source code from github and compile it locally.

#### Configure Your Server

It's recommended that you run a validator node on Linux Server.

**Recommended Configurations:**

1. 2 CPU
2. Memory: 6GB
3. Disk: 256GB SSD
4. OS: Ubuntu 16.04 LTS
5. Bandwidth: 20Mbps
6. Allow all incoming connections on TCP port 26656 and 26657

#### Install Go

::: tip
**Go 1.12.1+** is required for the IRIShub.
:::

Install `go` by following the [official docs](https://golang.org/doc/install).
 
::: tip
Quick install for macOS

Step 1: install `brew`

/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

Step 2: install `go` with `brew`

brew install go
:::

Verify that Go has been installed successfully：
```bash
go verison
```

Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
source ~/.bash_profile
```

#### Compile Source Code


- Get the code and compile Iris

After setting up Go correctly, you should be able to compile and run `iris`.
Make sure that your server can access to google.com because our project depends on some libraries provided by google.

* To compile for `testnet`:
Please checkout the latest version，refer to：https://github.com/irisnet/irishub/releases/latest
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

* To compile for `betanet`:
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

If your environment variables have set up correctly, you should not get any errors by running the above commands.
Now check your `iris` version.

```
$ iris version
<latest_iris_version>
    
$ iriscli version
<latest_iris_version>
```
