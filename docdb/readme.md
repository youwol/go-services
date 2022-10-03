# Youwol Document Database service

## How to build & install ?
* Make sure you have golang installed on your machine, and GOPATH pointing to your src folder containing this repo
* Download the dependecies/packages :
> go mod vendor
* Builds the executable in the current folder :
> go build ./cmd/docdb-server
* Builds and installs the executable in GOPATH/bin :
> go install ./cmd/docdb-server

## How to work with the API ?
* Download and install go-swagger from https://github.com/go-swagger/go-swagger/releases
* Edit api/docdb-api.yaml in http://editor.swagger.io
* Re-gen the code (models and endpoints)
* Go to this folder and then use the command:
> swagger generate server -P models.Principal /f api/docdb-api.yaml --name docdb

* How to test locally ?
1. Start scylladb and redis instances
> docker run --name scylla --rm -p 9042:9042 -d scylladb/scylla
> docker run --name redis -p 6379:6379 --rm -d redis
2. Start the docdb server
> docdb-server --port 8080
3. Then launch postman, and import the test collection api/DocDB_Tests.postman_collection.json and the local environment api/Local.postman_environment.json
The test collection can be played

## TODO
- Unit tests, including benchmarks
- Create index API
- Code quality tools automation
- GraphQL interface?
- Rationalize API responses
- Capability to stream select queries (if the number of rows is huge)