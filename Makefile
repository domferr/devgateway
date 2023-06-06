EXECUTABLE=bin/devgateway

build: build-linux

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${EXECUTABLE} ./src/main.go

run:
	./${EXECUTABLE}

start-postgresql:
	./${EXECUTABLE} -service postgresql

start-mysql:
	./${EXECUTABLE} -service mysql

start-mongodb:
	./${EXECUTABLE} -service mongodb

start-oracledb:
	./${EXECUTABLE} -service oracledb

start-harperdb:
	./${EXECUTABLE} -service harperdb

start-sqlserver:
	./${EXECUTABLE} -service sqlserver
