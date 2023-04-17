#
# Build image: docker build -t irisnet/irishub:v2.0.0 --build-arg EVM_CHAIN_ID=6688 .
#
FROM golang:1.18-alpine3.16 as builder

ARG EVM_CHAIN_ID

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash linux-headers eudev-dev

WORKDIR /irishub

# Add source files
COPY . .

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

RUN EVM_CHAIN_ID=$EVM_CHAIN_ID make build

# ----------------------------

FROM alpine:3.16

# p2p port
EXPOSE 26656
# rpc port
EXPOSE 26657
# metrics port
EXPOSE 26660

COPY --from=builder /irishub/build/ /usr/local/bin/