FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/typositoire/go-vln

COPY . .

ENV GO111MODULE=on

RUN go mod vendor && \
    go build -o /go/bin/go-vln

FROM alpine

COPY --from=builder /go/bin/go-vln /go/bin/go-vln

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && update-ca-certificates

ENTRYPOINT ["/go/bin/go-vln"]