FROM golang:alpine
ENV GOPATH /
ADD ./src/torpedobot /src/torpedobot
RUN apk update
RUN apk add git
RUN go get -v -d torpedobot
RUN go build -o /src/torpedo torpedobot
EXPOSE 3978
EXPOSE 3979
WORKDIR /src
ENTRYPOINT ./torpedo
