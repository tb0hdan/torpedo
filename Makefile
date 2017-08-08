.PHONY: build

ARCHITECTURE ?=
PLATFORM ?=
ARCHITECTURES = 386 amd64
PLATFORMS = darwin linux windows
GOPATH = $(shell pwd)
PKGNAME = "torpedobot"
DEST = $(PKGNAME)
BUILD = $(shell git rev-parse HEAD)
BDATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_UTC')
VERSION = $(shell cat ./VERSION)
BUILD_CMD = go build

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
	@go get -v -d $(PKGNAME)

report_deps:
	@go get -u -v github.com/wgliang/goreporter

build:  deps build_only

build_only:
	@$(BUILD_CMD) -ldflags "-X main.BUILD=$(BUILD) -X main.BUILD_DATE=$(BDATE) -X main.VERSION=$(VERSION)" -o bin/$(DEST) $(PKGNAME)


clean:
	@rm -rf bin/ build/ pkg/

coverage:
	@go test -bench=. -benchmem -race -cover torpedobot
 
coverage_html:	deps
	@go test -bench=. -benchmem -race -coverprofile=build/c.out $(PKGNAME)
	@go tool cover -html=build/c.out -o build/coverage.html
	@sleep 3; open http://localhost:8000/coverage.html
	@go run tools/fileserver.go -listen localhost:8000 -directory ./build

report:	clean deps report_deps
	@bin/goreporter -p ./src/torpedobot -r build/ -t src/github.com/wgliang/goreporter/templates/template.html

dockerimage:
	@cp -r /usr/local/etc/openssl ./ssl
	@docker build -t tb0hdan/torpedo .

dockerrun:
	@docker run --env-file ./token.sh tb0hdan/torpedo

release_binaries: deps
	@for platform in $(PLATFORMS); do for architecture in $(ARCHITECTURES); do echo "Building $(DEST)-$$platform-$$architecture"; make build_only PLATFORM=$$platform ARCHITECTURE=$$architecture; done; done
