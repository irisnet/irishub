# Simple usage with a mounted data directory:
# > docker build -t iris .
# > docker run -it --rm -v "/mnt/volumes/pangu:/iris" iris init [address]--home=/iris
# > docker run -it --rm -v "/mnt/volumes/pangu:/iris" iris start --home=/iris

FROM alpine:edge

ADD ./build/ /usr/local/bin/

ENV DATA_ROOT /iris

# Set user right away for determinism
RUN addgroup gaiauser && \
    adduser -S -G gaiauser gaiauser

# Create directory for persistence and give our user ownership
RUN mkdir -p $DATA_ROOT && \
    chown -R gaiauser:gaiauser $DATA_ROOT

VOLUME $DATA_ROOT

# p2p port
EXPOSE 46656
# rpc port
EXPOSE 46657

WORKDIR /bianjie/

ENTRYPOINT ["iris"]