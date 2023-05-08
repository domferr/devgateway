EXECUTABLE=bin/devgateway

build: build-linux

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${EXECUTABLE} ./src/main.go

run:
	./${EXECUTABLE}
