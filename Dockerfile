
# Simple usage with a mounted data directory:
# > docker build -t irishub .
# > docker run -v $HOME/.iris:/root/.iris iris init
# > docker run -v $HOME/.iris:/root/.iris iris start

FROM alpine:edge

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up GOPATH & PATH

ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/irisnet
ENV REPO_PATH    $BASE_PATH/irishub
ENV PATH         $GOPATH/bin:$PATH

# p2p port
EXPOSE 46656
# rpc port
EXPOSE 46657

# Add source files
COPY . $REPO_PATH/

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN cd $REPO_PATH && \
    apk add --no-cache $PACKAGES && \
    go get github.com/golang/dep/cmd/dep && \
    make get_vendor_deps && \
    make build_linux && \
    cp build/* /usr/local/bin/ && \
    cd / && \
    apk del $PACKAGES && \
    rm -rf $GOPATH/ && \
    rm -rf /root/.cache/

