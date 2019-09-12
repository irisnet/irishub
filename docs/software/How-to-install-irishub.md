# How to install `iris` 

## Latest Version

The Latest version of IRIShub is [v0.15.1](https://github.com/irisnet/irishub/releases/latest)

::: tip
Please replace <latest_iris_version> below with v0.15.1
:::

## Configure Your Server

It's recommended that you run irishub nodes on Linux Server.

**Recommended Configurations:**

- 2 CPU
- Memory: 6GB
- Disk: 256GB SSD
- OS: Ubuntu 16.04 LTS +
- Bandwidth: 20Mbps
- Allow all incoming connections on TCP port 26656 and 26657

## Install

### Install `go`

::: tip
**Go 1.12.5+** is required for the IRIShub.
:::

Install `go` by following the [official docs](https://golang.org/doc/install).
 
Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN" >> ~/.bashrc
source ~/.bashrc
```

Verify that Go has been installed successfully: 
```bash
go version
```

### Install `iris`

After setting up Go correctly, you should be able to compile and run `iris`.

Make sure that your server can access to google.com because our project depends on some libraries provided by google. (If you are not able to access google.com, you can also try to add a proxy: `export GOPROXY=https://goproxy.io`)

```bash
git clone --branch <latest_iris_version> https://github.com/irisnet/irishub
cd irishub
# source scripts/setTestEnv.sh # used for installing the testnet version
make get_tools install
```

If your environment variables have set up correctly, you should not get any errors by running the above commands.
Now check your `iris` version.

```
$ iris version
<latest_iris_version>
    
$ iriscli version
<latest_iris_version>
```
