all: build

build:
	GOOS=linux GOARCH=amd64 go build -v -o bin/smally-go ./...

format:
	go fmt ./...

test check:
	go test -race --coverprofile=coverage.coverprofile --covermode=atomic ./...

lint:
	golint ./...
	go vet ./...

test-unit:
	go test -v ./...

run-local:
	go build -v -o bin/smally-go ./...
	./bin/smally-go/server -http.port=":8080" -redis.host="localhost" -redis.port="6379"