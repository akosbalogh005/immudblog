SKIPPED_UNITTESTS_ON_DIRS=mocks|docs|samples

build: dep
	go build -o bin/immudblog-server cmd/immudblog-server/main.go
	go build -o bin/immudblog-cli cmd/immudblog-cli/main.go

dep:
	go mod tidy

swaggo-gen:
	swag init -d cmd/immudblog-server,./restapi,./model

clean:
	echo ${VERSION}
	rm -rf bin
	
generate:
	find -name mocks -exec rm -r {} +
	# mockery is needed: 'go get github.com/vektra/mockery/v2/.../'
	go generate ./...

test: dep
	go test `go list ./... | egrep -v "${SKIPPED_UNITTESTS_ON_DIRS}"` -coverprofile coverage.txt -covermode count && go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt -o coverage-report.html

testwithresults: test
	go tool cover -html=coverage.txt	