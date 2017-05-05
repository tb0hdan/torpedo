GOPATH=$(shell pwd)

all: build

deps:
	@go get -v -d torpedobot

build:  deps
	@go build -o bin/torpedobot torpedobot

test:
	go test -bench=. -benchmem -race -cover torpedobot

clean:
	rm -f bin/torpedobot
