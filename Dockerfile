#
# Build image: docker build -t irisnet/irishub:v2.1.0 --build-arg EVM_CHAIN_ID=6688 .
# Base images can be overridden when Docker Hub is unavailable.
#
ARG BUILDER_IMAGE=golang:1.24.9-alpine
ARG RUNTIME_IMAGE=alpine:3.18

FROM ${BUILDER_IMAGE} AS builder

ARG EVM_CHAIN_ID

# Set up dependencies
ENV PACKAGES="make gcc git libc-dev bash linux-headers eudev-dev build-base"

WORKDIR /irishub

# Add source files
COPY . .

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

RUN EVM_CHAIN_ID=$EVM_CHAIN_ID make build

# ----------------------------

FROM ${RUNTIME_IMAGE}

# p2p port
EXPOSE 26656
# rpc port
EXPOSE 26657
# metrics port
EXPOSE 26660

COPY --from=builder /irishub/build/ /usr/local/bin/
