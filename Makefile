GOFILES     = $(wildcard *.go)
GONAME      = $(shell basename "$(PWD)")
PACKAGES	= $(shell go list ./...)

#tests:
#	@echo "Launching Tests in Docker Compose"
#	mkdir -p ./testdata
#	test -f ./testdata/testkey || ssh-keygen -b 2048 -t rsa -f ./testdata/testkey -q -N ""
#	docker-compose -f dev-compose.yml up --build tests

#clean:
#	@echo "Cleaning up build junk"
#	-docker-compose -f dev-compose.yml down
#	-rm -rf ./testdata


build:
	@echo "Building $(GOFILES) to ./bin"
	GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

install:
	@echo "Installing from source"
	WORKSPACE=$(pwd -P)
	go install

fmt:
	go fmt $(PACKAGES)

lint:
	golint ./...