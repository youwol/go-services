# Libraries for golang services

## How-to test ?
- install gotestsum (a testing utility)
```
go get gotest.tools/gotestsum
```
- run tests
```
cd test
gotestsum --format testname --junitfile ../junit.xml
```
