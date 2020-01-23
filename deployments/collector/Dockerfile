FROM ubuntu:18.04

RUN apt-get update && \
  apt-get -y install \
  vim net-tools iputils-ping ca-certificates

# Load Go binary executable
COPY collector /
COPY locations.json /

# Launch it at startup
ENTRYPOINT /collector