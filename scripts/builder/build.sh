#!/bin/bash

function display_usage {
   echo -e ""
   echo -e "Builds an image to run River Life builds"
   echo -e ""
   echo -e "Usage: $FILENAME <image_name> <version>"
   echo -e ""
}


IMAGE_NAME=$1
IMAGE_VERSION=$2

if [ $# -ne 2 ];then
   echo "Invalid parms: see usage"
   display_usage
   exit -1
fi

echo "image_name=$IMAGE_NAME"
echo "version=$IMAGE_VERSION"

# Location of the source code to map into the container
# (sourced from a local environment variable)
export RLCODE="${RLCODE:-/root/repo}"

# First, go build the container on the local host system
echo -e "Building ${IMAGE_NAME} image"
docker build \
   --tag="${IMAGE_NAME}:${IMAGE_VERSION}" \
   ./builder
exit $?
