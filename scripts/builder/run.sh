#!/bin/bash

function display_usage {
   echo -e ""
   echo -e "Runs the River Life build container"
   echo -e ""
   echo -e "Usage: $FILENAME <container_name> <image_version>"
   echo -e ""
}

function stop_container {
   echo -e "Stopping an already running instance of $1"
   docker rm $1 -f
}

CONTAINER_NAME=$1
IMAGE_VERSION=$2

if [ $# -ne 2 ];then
   echo "Invalid parms: see usage"
   display_usage
   exit -1
fi

echo "container_name=$CONTAINER_NAME"
echo "image_version=$IMAGE_VERSION"

# Location of the source code to map into the container
# (sourced from a local environment variable)
export RLCODE="${RLCODE:-/root/repo}"

# Second make sure the same container isn't already running
CONTAINER_RUNNING=$(docker ps -a | grep -w ${CONTAINER_NAME})
if [ -n "${CONTAINER_RUNNING}" ]; then
   stop_container $CONTAINER_NAME
fi

echo -e "\n\nStarting $CONTAINER_NAME container..."
docker run -d -it \
   --name ${CONTAINER_NAME} \
   --hostname ${CONTAINER_NAME} \
   --restart unless-stopped \
   -v /$RLCODE:/root/rlcode/src/riverlife \
   -w "//root/rlcode" \
   ${CONTAINER_NAME}:${IMAGE_VERSION}
exit $?