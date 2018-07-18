#!/bin/bash
docker run -d -p 46656:46656 -p 46657:46657 --name irishub-sandbox -it irishub:develop sh -c '
    apk add --no-cache curl expect && \
    curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/start.sh -o start.sh && \
    sh start.sh'