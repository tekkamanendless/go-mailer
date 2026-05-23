all: cicd-build

RACE ?= 0
export CGO_ENABLED ?= 0
ifeq ($(RACE), 1)
	GO_RACE := -race
	export CGO_ENABLED := 1
else
	GO_RACE :=
endif

ALL_GO_FILES := $(shell find ./ -name '*.go')

current_dir = $(shell pwd)

.PHONY: clean
clean:
	go clean

.PHONY: upgrade-deps
upgrade-deps:
	go get -u ./...
	go get -u all
	go get -u
	go mod tidy

.PHONY: test
test:
	go vet ./...
	go test -cover $(GO_RACE) -parallel 10 ./...

.PHONY: staticcheck
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@v0.6.1
	staticcheck ./...

.PHONY: format
format:
	go fmt ./...

#
# CI/CD targets
#

.PHONY: cicd-build
cicd-build:
	go build ./...

