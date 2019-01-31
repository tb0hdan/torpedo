.PHONY: build

ARCHITECTURE ?=
PLATFORM ?=
GO = "go"
ARCHITECTURES = 386 amd64
PLATFORMS = darwin linux windows
GOPATH = $(shell pwd)
PKGNAME = "torpedobot"
DEST = $(PKGNAME)
BUILD = $(shell git rev-parse HEAD)
BDATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_UTC')
GO_VERSION = $(shell $(GO) version|awk '{print $$3}')
VERSION = $(shell cat ./VERSION)
BUILD_CMD = $(GO) build

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
	@mkdir -p bin/ build/ pkg/
	@$(GO) get -v -d $(PKGNAME)

report_deps:
	@$(GO) get -u -v github.com/wgliang/goreporter

build:  deps build_only

build_only:
	@$(BUILD_CMD) -v -x -ldflags "-X main.BUILD=$(BUILD) -X main.BUILD_DATE=$(BDATE) -X main.GO_VERSION=$(GO_VERSION) -X main.VERSION=$(VERSION)" -o bin/$(DEST) $(PKGNAME)


clean:
	@rm -rf bin/ build/ pkg/

coverage:
	@$(GO) test -bench=. -benchmem -race -cover $(PKGNAME)
 
coverage_html:	deps
	@$(GO) test -bench=. -benchmem -race -coverprofile=build/c.out $(PKGNAME)
	@$(GO) tool cover -html=build/c.out -o build/coverage.html
	@sleep 3; open http://localhost:8000/coverage.html
	@$(GO) run tools/fileserver.go -listen localhost:8000 -directory ./build

report:	clean deps report_deps
	@bin/goreporter -p ./src/$(PKGNAME) -r build/ -t src/github.com/wgliang/goreporter/templates/template.html

codecov:
	@$(GO) test -race -coverprofile=coverage.txt -covermode=atomic $(PKGNAME)

dockerimage:
	@cp -r /usr/local/etc/openssl ./ssl
	@docker build -t tb0hdan/torpedo .

dockerrun:
	@docker run --env-file ./token.sh tb0hdan/torpedo

release_binaries: deps
	@for platform in $(PLATFORMS); do for architecture in $(ARCHITECTURES); do echo "Building $(DEST)-$$platform-$$architecture"; make build_only PLATFORM=$$platform ARCHITECTURE=$$architecture; done; done

race:
	@$(GO) test -race $(PKGNAME)

trace:
	@$(GO) test -bench=. -trace trace.out $(PKGNAME)
	@$(GO) tool trace trace.out

lint:
	@golangci-lint run ./src/torpedobot

tag:
	@git tag -a v$(VERSION) -m v$(VERSION)
	@git push --tags
