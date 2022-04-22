#!/bin/bash
docker-compose -f ../../docker-compose.test.yml up -d test-db 2>&1
export ITLAB_REPORTS_MONGO_URI=mongodb://root:root@localhost:27018/itlab-reports?authSource=admin
export ITLAB_REPORTS_MONGO_TEST_URI=mongodb://root:root@localhost:27018/itlab-reports-test?authSource=admin

go install github.com/jstemmer/go-junit-report 2>&1
go install github.com/axw/gocov/gocov 2>&1
go install github.com/AlekSi/gocov-xml 2>&1
# Run Go tests and turn output into JUnit test result format
touch TestCoverage.txt TestReport.xml
sleep 20
go test ./... -p 2 -v -coverprofile=TestCoverage.txt -covermode count 2>&1 | $HOME/go/bin/go-junit-report > TestReport.xml
rc=${PIPESTATUS[0]} # Get result code of `go test`

if [ $rc -ne 0 ]; then
    # Let script fail by writing to stderr
    >&2 echo "Tests failed"
fi
