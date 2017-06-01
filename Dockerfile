# build stage
FROM golang:alpine AS build-env
ENV GOPATH /
WORKDIR /src
ADD ./ssl /etc/ssl
ADD ./src/torpedobot /src/torpedobot
RUN apk update
RUN apk add git
RUN go get -v -d torpedobot
RUN go build -o /src/torpedo torpedobot

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /etc/ssl /etc/ssl
COPY --from=build-env /src/torpedo /app/
EXPOSE 3978
EXPOSE 3979
EXPOSE 3980
EXPOSE 3981
ENTRYPOINT ./torpedo
