.PHONY: all
## Builds and verifies the project.
all: build test lint

.PHONY: pipeline
## Runs the same thing as the pipeline.
pipeline: all

.PHONY: build
## Builds (compiles) the project.
build: deps
	go build ./...

.PHONY: deps
## Gets the dependencies.
deps:
	go get -v -t -d ./...

.PHNOY: test
## Runs the tests.
test:
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
