FROM ubuntu:18.04

RUN apt-get update && \
  apt-get -y install \
  vim net-tools iputils-ping openssl

# Load Go binary executable
COPY  apisvr /

# Launch it at startup
ENTRYPOINT /apisvr