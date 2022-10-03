# Youwol Storage service

## How to build & install ?
Make sure you have golang installed on your machine, and GOPATH pointing to your src folder containing this repo (e.g. ~myprofile/src/platform, GOPATH is set to ~myprofile/src)

* Download the dependecies/packages
> go mod vendor
* Build the executable in the current folder
> go build ./cmd/storage-server
* Builds and installs the executable in GOPATH/bin
> go install ./cmd/storage-server

## How to work with the API ?
1. Download and install go-swagger from https://github.com/go-swagger/go-swagger/releases
2. Edit api/storage-api.yaml in http://editor.swagger.io
3. Re-gen the code (models and endpoints)
4. Go to this folder and then use the command:
> swagger generate server -P models.Principal /f api/storage-api.yaml

## How to test locally ?
1. Start a minio instance
> docker run --name minio --rm -p 9000:9000 -v D:\home\tmp\minio:/data -d minio/minio server /data
2. Start a redis instance
> docker run --name redis -p 6379:6379 --rm -d redis
3. Start the storage server
> storage-server --port 8090
4. Then launch postman, and import the test collection ../test/postman/suites/Storage_Tests.postman_collection.json and the local environment api/Local.postman_environment.json

## TODO
* Data authorization
* Unit tests, including benchmarks
* Create notification API
* Code quality tools automation
* Security, API Keys...
* (optional) Test suite: make file upload dynamic with variables and using postman working directory
