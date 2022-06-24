# ITLab-Reports
Service for storing work reports

Status | master | develop
---|---|---
build | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=develop)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
test | [![master tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/master?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![develop tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/develop?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
## Requirements
- Go 1.18+ || Docker
## Configuration

### from enviroment
```bash
# url to jwks.json
ITLAB_REPORTS_AUTH_KEY_URL=https://examplesite/files/jwks.json

# audince for JWT
# in default itlab
ITLAB_REPORTS_AUTH_AUDIENCE=audience

# issuer for JWT  
ITLAB_REPORTS_AUTH_ISSUER=https://exampleissuersite.com

# required scope for JWT
ITLAB_REPORTS_AUTH_SCOPE=my_scope

# url to database
ITLAB_REPORTS_MONGO_URI=mongodb://user:password@localhost:27017

# url to test database for launching tests
ITLAB_REPORTS_MONGO_URI=mongodb://user:password@localhost:27017

# app running port 
ITLAB_REPORTS_APP_PORT=8080

# grpc application port
ITLAB_REPORTS_APP_GRPC_PORT=8081

# testMode=true disables jwt validation
ITLAB_REPORTS_APP_TEST_MODE=false

# user role
# in deafult user
ITLAB_REPORTS_ROLE_USER=user

# admin role
# in deafult reports.admin
ITLAB_REPORTS_ROLE_ADMIN=admin

# superadmin role
# in deafult admin
ITLAB_REPORTS_ROLE_SUPER_ADMIN=superadmin

# base url for external api
# in default https://dev.manage.rtuitlab.dev
ITLAB_REPORTS_APP_ITLAB_URL=https://some_api_baseurl.com
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
1. Build binary file
    ```bash
    cd src/ITLabReports/api
    go build -o main
    ```
1. Launch file
    ```bash
    ./main
    ```
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
    1. Build binary file
        ```bash
        cd src/ITLabReports/api
        go build -o main
        ```
    1. Launch file
        ```bash
        ./main
        ```
1. Run tests with Karate
    1. Install JDK 1.8 (ensure JAVA_HOME env variable exist)
    1. Install Maven
    1. Run tests
        ```bash
        cd tests/e2e/Karate
        mvn clean install '-Dmaven.test.skip=true'
        mvn test '-Dtest=TestParallel.java'
        ```
1. Run tests with Testmace
    1. Install testmace dependency
        ```bash
        npm install --global @testmace/cli@1.3.1
        ```
    1. Run tests
        ```bash
        cd tests/e2e/TestMace
        testmace-cli ./Project --reporter=junit -e localEnv -o tests-out
        ```
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
## Documantation
docs available on /api/reports/swagger/

## Work with test mode


in testmode use test auth middleware
to generate token for it use https://jwt.io
use secret "test" and put into sub id of your test user
and in audiance field that configure with env put your role from env role for example for default test env use next payload

```json
{
  "aud": "audiance",
  "iss": "https://example.com",
  "iat": 1516239022,
  "exp": 1505467756869,
  "sub": "user_id",
  "itlab": [
    "user"
  ]
}
```
