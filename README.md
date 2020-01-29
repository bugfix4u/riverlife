# **River Life**

## Description
This code is part of the River Life project I am developing, to build an Amazon Alexa skill to report river hydrolic observations from NOAA collection sites. This app is currently made up of three micro services: a collection service, a postgres database, and an API server. The collection service is a multi-threaded service that pulls down and parses RSS feeds for each state, and from those RSS feeds it collects information from over 11,000+ collection sites around the United States. This data is then parsed and put in a Postgres DB for use by the REST API server. This is still a very much a "work in progress", so be sure to check back often.

## Requirements
To build or run this code for yourself, you will need to have a linux machine with docker installed, and an internet connection. The build scripts are currently built in bash, and all of the compiling and running of code is done through docker containers. See the "Building and Running the Code" section for more details

## Building and Running the Code
### Build Script
In the root of the Riverlife directory is a `build.sh` script which can be used to compile, build docker containers, run the docker containers, package the containers for export, and install the container into docker. All files resulting from the compile and build commands are places in the `./builds` directory.

### Compiling the Riverlife Code
To compile the Riverlife code you can run `sudo ./build.sh compile`. This command will build a container image to compile the Riverlife code and run the container in docker. The first time you run the compile command, it will take several minutes to build the compiler image and load it. Subsequent compiles will be much fast, taking only a few seconds, since the build container is left running until manually stopped or the version number changes.

### Building Docker Images
To build the docker images you can run `sudo ./build.sh build`. This command will copy over any needed Dockerfiles, docker-compose, configuration, etc. It will then download base images from the Docker Hub, and install any needed application, updates, code, and properties files. Lastly it will save a copy of the containewr image in the `./builds` directory for packaging at a later stage.

### Running the Riverlife Microservices
To start the RiverLife microservices  you can run `sudo ./build.sh run`. This command will stop and remove any current RiverLife containers that are running, and will then create new containers from the images in the local docker repo.

### Packaging the Riverlife Microservices
To package the Riverlife microservices you can run `sudo ./build.sh package`. This will package up all of the container images, docker-compose files, env files, scripts, etc into a compressed tar image. As part of this tar image the `rl_run` and `rl_install` scripts will be included. This will make it easier to move to another docker server or VM for running.

### Installing the Riverlife Microservices
Instructions for this will be coming soon, after more testing is done.

## Watching the Microservice Logs
You can tail the logs for the rl-apisvr and rl-collector services by using the following command:
```bash
sudo docker exec -it <container name> tail -f /riverlife.log
```

## ToDo
This is the current ToDo list in no particular order or priority

- Implement HTTP HEAD check to check if data has changed before downloading XML
- Add a Redis service to handle some caching for the collection service
- Create additional REST API's for seaching the data by different parameters
- Add paging support to REST API's that return lists of information
- ...
