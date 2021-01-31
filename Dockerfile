FROM golang:1.16rc1 AS builder

WORKDIR /go/src/
COPY hello.go .

RUN go build -o hello hello.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/hello .
CMD ["./hello"]