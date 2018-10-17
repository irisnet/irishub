# Protobuf

##Install g++
g++ version should be no less than 7

For ubuntu operator system, please follow these steps:
```
sudo add-apt-repository ppa:ubuntu-toolchain-r/test
sudo apt update
sudo apt install g++-7 -y
```

##Install Protocol Buffers Library

Please follow this [document](https://github.com/protocolbuffers/protobuf/blob/master/src/README.md) to install protocol buffers library.

##Build
Build cpp to library:
```
make build
```
Clean the library
```
make clean
```