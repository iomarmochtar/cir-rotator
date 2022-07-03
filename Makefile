NAME           ?= "cir-rotator"
LINTER_VERSION  = 1.46.0
LINTER_BIN      = ./bin/golangci-lint
COVERAGE_OUT    = coverage.txt

.PHONY: setup-linter
export LINTER_VERSION
setup-linter:
	[ ! -f ./myfile ] && curl -sfL https://raw.githubusercontent.com/bukalapak/toolkit-installer/master/golangci-lint.sh | sh

.PHONY: lint
lint: setup-linter
	${LINTER_BIN} run -v

.PHONY: test
test:
	go test -v -cover -coverprofile=${COVERAGE_OUT} -covermode=atomic ./...

.PHONY: coverage-html
coverage-html: test
	go tool cover -html=${COVERAGE_OUT}

.PHONY: test-all
test-all: test lint

.PHONY: compile
compile:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o dist/${NAME} main.go
