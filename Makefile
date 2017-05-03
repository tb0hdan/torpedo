GOPATH=$(shell pwd)

all: build

build:
	go get -v -d torpedobot
	go build -o bin/torpedobot torpedobot

test:
	go test -bench=. -benchmem -race -cover torpedobot

clean:
	rm -f bin/torpedobot
