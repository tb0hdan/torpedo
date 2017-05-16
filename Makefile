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

dockerimage:
	@docker build -t tb0hdan/torpedobot .

dockerrun:
	@docker run --env-file ./token.sh tb0hdan/torpedobot
