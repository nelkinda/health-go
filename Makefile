.PHONY: all
## Builds and verifies the project.
all: build test

.PHONY: pipeline
## Runs the same thing as the pipeline.
pipeline: all

.PHONY: build
## Builds (compiles) the project.
build: deps generate
	go build ./...

.PHONY: deps
## Gets the dependencies.
deps: generate
	go get -v -t -d ./...

.PHONY: generate
## Generates source code.
generate:
	go generate ./...

.PHNOY: test
## Runs the tests.
test:
	go test -test.cover .
	go test ./...

.PHONY: lint
## Verifies the source code using golint.
lint:
	golint -set_exit_status ./...

-include .makehelp/include/makehelp/Help.mk

ifeq (help, $(filter help,$(MAKECMDGOALS)))
.makehelp/include/makehelp/Help.mk:
	git clone --depth=1 https://github.com/christianhujer/makehelp.git .makehelp
endif
