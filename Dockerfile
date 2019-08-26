FROM golang:buster

WORKDIR /go/src/api
COPY . .

RUN go get -d -v ./...
RUN go build -v .

FROM debian:buster
COPY --from=0 /go/src/api/api /usr/bin/api

ENV LISTEN_ADDR="localhost:9090"
ENV POSTGRES_HOST=""
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=""
ENV POSTGRES_PASSWORD=""
ENV POSTGRES_DATABASE=""

ENTRYPOINT ["/usr/bin/api"]