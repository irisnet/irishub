#
# Build image: docker build -t irisnet/irishub .
#
FROM golang:1.15.5-alpine3.12 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash linux-headers eudev-dev

WORKDIR /irishub

# Add source files
COPY . .

# Install minimum necessary dependencies, run unit tests
RUN apk add --no-cache $PACKAGES && make test-unit

RUN make build

# ----------------------------

FROM alpine:3.12

# p2p port
EXPOSE 26656
# rpc port
EXPOSE 26657
# metrics port
EXPOSE 26660

COPY --from=builder /irishub/build/ /usr/local/bin/