.EXPORT_ALL_VARIABLES:
.PHONY:
.SILENT:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := start

include .env

dev:
	air

start:
	GIN_MODE=release go run cmd/main.go

#testing variables
export TEST_DB_URI=mongo://$(MONGODB_USER):$(MONGODB_PASSWORD)@$(MONGODB_HOST):$(MONGODB_PORT)
export TEST_DB_NAME=test
export TEST_CONTAINER_NAME=test_db

tests:
	go test -v ./test/...

lint:
	golangci-lint run

swagger:
	swag init -g cmd/main.go

env:
	echo $$TEST_DB_URI