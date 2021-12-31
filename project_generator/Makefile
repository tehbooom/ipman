BINARY_NAME=projectGen

build:
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux *.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows *.go

run: 
	./bin/${BINARY_NAME}-linux 

all: build run

clean:
	go clean
	rm ${BINARY_NAME}-main-linux 
	rm ${BINARY_NAME}-main-windows

lint:
	golangci-lint run
	
test:

dep:
	go mod download