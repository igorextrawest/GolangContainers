FROM golang:1.11 as builder

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/igorextrawest/GolangContainers

ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

RUN dep ensure --vendor-only

ADD src src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app src/main.go

RUN go test -v ./...

FROM alpine:3.7

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root

COPY --from=builder /go/src/github.com/igorextrawest/GolangContainers/app .

CMD ["./app"]