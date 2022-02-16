# ITLab-Reports
Service for storing work reports

Status | master | develop
---|---|---
build | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=develop)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
test | [![master tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/master?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![develop tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/develop?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
## Requirements
- Go 1.13.8+ || Docker
## Configuration
File `src/api/auth_config.json` must contain next content:
```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", // url to jwks.json       | env: ITLAB_REPORTS_AUTH_KEY_URL
    "audience": "example_audience",                  // audince for JWT        | env: ITLAB_REPORTS_AUTH_AUDIENCE
    "issuer" : "https://exampleissuersite.com",      // issuer for JWT         | env: ITLAB_REPORTS_AUTH_ISSUER
    "scope" : "my_scope",                            // required scope for JWT | env: ITLAB_REPORTS_AUTH_SCOPE
  }
}
```  
File `src/api/config.json` must contain next content:
```js
{
  "DbOptions": {
    "uri": "mongodb://user:password@localhost:27017", // url to database          | env: ITLAB_REPORTS_MONGO_URI
    "dbName" : "ITLabReports",                        // database name            | env: ITLAB_REPORTS_MONGO_DB_NAME
    "collectionName" : "reports",                     // databsae collection name | env: ITLAB_REPORTS_MONGO_DB_COLLECTION_NAME
  },
  "AppOptions": {
    "appPort": "8080", // app running port                      | env: ITLAB_REPORTS_APP_PORT
    "testMode": false  // testMode=true disables jwt validation | env: ITLAB_REPORTS_APP_TEST_MODE
  }
}
```
## Run locally
1. Complete [configuration](#configuration)
### Via Docker
1. Build image
    ```bash
    docker build -t rtuitlab_reports-back -f Dockerfile .
    ```
1. Start container
    ```bash
    docker run -d -p 8080:8080 rtuitlab_reports-back
    ```
### Via Golang
1. TODO: steps to build ITLab-Reports via Golang
## Run tests
### Via Docker
1. Run app
    ```bash
    docker compose -f docker-compose.test.yml up -d test-api
    ```
1. Run tests
    1. With Karate
        ```bash
        docker compose -f docker-compose.test.yml up karate
        ```
    1. With Testmace
        ```bash
        docker compose -f docker-compose.test.yml up testmace
        ```
1. Result are stored in `tests/e2e/out-karate` and `tests/e2e/out-testmace` folders
### Manually
1. Run app
    ```bash
    # TODO: script to run ITLab-Reports via Golang
    ```
1. Run tests with Karate
    1. Install JDK 1.8 (ensure JAVA_HOME env variable exist)
    1. Install Maven
    1. Run tests
        ```bash
        cd tests\e2e\Karate
        mvn clean install '-Dmaven.test.skip=true'
        mvn test '-Dtest=TestParallel.java'
        ```
1. Run tests with Testmace
    1. TODO: steps to run Testmace
## Get Docker stack config
1. Generate docker stack
    ```bash
    docker compose -f docker-compose.yml -f docker-compose.prod.yml config
    ```
## Run in [ITLab](https://github.com/RTUITLab/ITLab) local stack
Read [ITLab README.md](https://github.com/RTUITLab/ITLab/blob/master/README.md) to generate self-signed certificate
## Requests
You can get Postman requests collection [here](https://www.getpostman.com/collections/4085657bcce140031d0c)
## DB Backup and Restore
To make a backup of a DB, open root folder where MongoDB is installed, open a command promt and type the command mongodump  
To restore the backup, open root folder where MongoDB is installed, open a command promt and type the command mongorestore  
(All DB paths are default)
