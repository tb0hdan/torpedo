.PHONY: build

ARCHITECTURE ?=
PLATFORM ?=
GO111MODULE = on
GO = GO111MODULE=$(GO111MODULE) go
ARCHITECTURES = 386 amd64
PLATFORMS = darwin linux windows
# Preserved for historical purposes
# GOPATH = $(shell pwd)
PKGNAME = "torpedobot"
PROJECT_URL = "https://github.com/tb0hdan/torpedo"
DEST = $(PKGNAME)
BUILD = $(shell git rev-parse HEAD)
BDATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_UTC')
GO_VERSION = $(shell $(GO) version|awk '{print $$3}')
VERSION = $(shell cat ./VERSION)
BUILD_CMD =  $(GO) build

ifneq ($(strip $(PLATFORM)),)
    BUILD_CMD := GOOS=$(PLATFORM) $(BUILD_CMD)
    DEST := $(DEST)-$(PLATFORM)
endif

ifneq ($(strip $(ARCHITECTURE)),)
    BUILD_CMD := GOARCH=$(ARCHITECTURE) $(BUILD_CMD)
    DEST :=  $(DEST)-$(ARCHITECTURE)
endif



all: build

deps:
	@mkdir -p bin/ build/
	@cd src/$(PKGNAME); $(GO) mod verify
	@cd src/$(PKGNAME); $(GO) mod download

report_deps:
	@cd src/$(PKGNAME); $(GO) get -u -v github.com/wgliang/goreporter

build:  deps lint build_only

build_only:
	@cd src/$(PKGNAME); $(BUILD_CMD) -v -x -ldflags "-X main.Build=$(BUILD) -X main.BuildDate=$(BDATE) -X main.GoVersion=$(GO_VERSION) -X main.Version=$(VERSION) -X main.ProjectURL=$(PROJECT_URL)" -o ../../bin/$(DEST) $(PKGNAME)

clean:
	@rm -rf bin/ build/

coverage:
	@cd src/$(PKGNAME); $(GO) test -bench=. -benchmem -race -cover $(PKGNAME)
 
coverage_html:	deps
	@cd src/$(PKGNAME); $(GO) test -bench=. -benchmem -race -coverprofile=../../build/c.out $(PKGNAME)
	@cd src/$(PKGNAME); $(GO) tool cover -html=../../build/c.out -o ../../build/coverage.html
	@sleep 3; open http://localhost:8000/coverage.html
	@$(GO) run tools/fileserver.go -listen localhost:8000 -directory ./build

report:	clean deps report_deps
	@bin/goreporter -p ./src/$(PKGNAME) -r build/ -t src/github.com/wgliang/goreporter/templates/template.html

codecov:
	@cd src/$(PKGNAME); $(GO) test -race -coverprofile=coverage.txt -covermode=atomic $(PKGNAME)

dockerimage:
	@cp -r /usr/local/etc/openssl ./ssl
	@docker build -t tb0hdan/torpedo .

dockerrun:
	@docker run --env-file ./token.sh tb0hdan/torpedo

release_binaries: deps
	@for platform in $(PLATFORMS); do for architecture in $(ARCHITECTURES); do echo "Building $(DEST)-$$platform-$$architecture"; make build_only PLATFORM=$$platform ARCHITECTURE=$$architecture; done; done

race:
	@cd src/$(PKGNAME); $(GO) test -race $(PKGNAME)

trace:
	@cd src/$(PKGNAME); $(GO) test -bench=. -trace trace.out $(PKGNAME)
	@cd src/$(PKGNAME); $(GO) tool trace trace.out

lint:
	@if [ "$(shell which golangci-lint)" == "" ]; then go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; fi
	@cd src/$(PKGNAME); GO111MODULE=$(GO111MODULE) golangci-lint run -v -c ../../.golangci.yml .

tag:
	@git tag -a v$(VERSION) -m v$(VERSION)
	@git push --tags

env:
	@env
