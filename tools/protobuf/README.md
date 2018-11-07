# Protobuf

##Construct build environment

1. Install g++
    g++ version should be no less than 7. For ubuntu operator system, please follow these steps:
    ```
    sudo add-apt-repository ppa:ubuntu-toolchain-r/test
    sudo apt update
    sudo apt install g++-7 -y
    ```

2. Install Protocol Buffers Library
    Please follow this [document](https://github.com/protocolbuffers/protobuf/blob/master/src/README.md) to install protocol buffers library.

3. Install libjsoncpp-dev
    ```
    sudo apt-get install libjsoncpp-dev
    sudo ln -s /usr/include/jsoncpp/json/ /usr/include/json
    ```

##Build
1. Build cpp to library:
    ```
    make build
    ```
2. Clean the library
    ```
    make clean
    ```