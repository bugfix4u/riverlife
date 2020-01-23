# **River Life**

## Description
This code is part of the River Life project I am developing, to build an Amazon Alexa skill to report river hydrolic observations from NOAA collection sites. This app is currently made up of three micro services: a collection service, a postgres database, and an API server. The collection service is multi-threaded service that pulls down and parses RSS feeds for each state, and from those RSS feeds is collects over 11,000+ collection sites around the United States. This data is then parsed and put in a Postgres DB for use by the REST API server. This is still a work in progress.

## ToDo
This is the current ToDo list in no particular order or priority

- Add logrus logging to API server
- Implement HTTP HEAD check to check if data has changed before downloading XML
- Add a Redis service to handle some caching for the collection service
- Create additional REST API's for seaching the data by different parameters
- Add paging support to REST API's that return lists of information
- Add configuration to overide defaults for logging, threads, polling times, etc
- ...
