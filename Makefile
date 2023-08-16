BINARY_NAME=quake

build:
	go build -o bin/$(BINARY_NAME) 

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe

run: build
	./bin/$(BINARY_NAME)

clean:
	go clean
	rm -rf ./bin
	rm -f ./coverage.out
	rm -f ./quake

vendor:
	go mod vendor

dep:
	go mod download

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-build:
	docker build -t quake:v0.0.1 .

docker-run: docker-build
	docker compose up -d

docker-stop:
	docker compose down