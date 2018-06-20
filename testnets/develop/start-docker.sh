#!/bin/bash
docker run -d -p 46656:46656 -p 46657:46657 -it irishub:develop sh -c '
    iris init --chain-id=fuxi-develop --name=init1 && \
    apk add --no-cache curl && \
    curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/genesis.json -o ~/.iris/config/genesis.json && \

    iris start'