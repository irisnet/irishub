#
# Build mainnet image: docker build -t irisnet/irishub .
# Build testnet image: docker build -t irisnet/irishub --build-arg NetworkType=testnet .
#
FROM golang:1.12.5-alpine3.9 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash linux-headers eudev-dev

WORKDIR /irishub

# Add source files
COPY . .

# Install minimum necessary dependencies, run unit tests
RUN apk add --no-cache $PACKAGES && make get_tools && make test_unit

# Initialize network type, could be override via docker build argument `--build-arg NetworkType=testnet`
ARG NetworkType=mainnet

RUN make build_cur

FROM alpine:3.9

# p2p port
EXPOSE 26656
# rpc port
EXPOSE 26657
# metrics port
EXPOSE 26660

COPY --from=builder /irishub/build/ /usr/local/bin/