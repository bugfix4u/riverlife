FROM ubuntu:18.04

RUN apt-get update && \
   apt-get -y install \
   wget vim git

# Create generic dir for installing packcages
RUN mkdir /install

# Install Golang v1.13.6
WORKDIR /install
RUN wget https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
RUN tar -xvf go1.13.6.linux-amd64.tar.gz
RUN mv go /usr/local

# Setup Go build environment
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/rlcode
ENV PATH=$PATH:$GOPATH:$GOPATH/bin:$GOROOT/bin
RUN mkdir -p $GOPATH/src/riverlife

# Run the build script
COPY gobuild.sh $GOPATH
RUN chmod 755 $GOPATH/gobuild.sh
# ENTRYPOINT ["/root/rlcode/gobuild.sh"]