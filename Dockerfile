# build stage
FROM golang:alpine AS build-env
ENV GOPATH /
WORKDIR /
ADD ./ssl /etc/ssl
ADD ./Makefile /
ADD ./VERSION /
ADD ./.git /
ADD ./src/torpedobot /src/torpedobot
RUN apk update
RUN apk add git
RUN apk add make
RUN apk add gcc
RUN apk add libc-dev
RUN make deps
RUN make build_only

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /etc/ssl /etc/ssl
COPY --from=build-env /bin/torpedobot /app/
EXPOSE 3978
EXPOSE 3979
EXPOSE 3980
EXPOSE 3981
ENTRYPOINT ./torpedobot
