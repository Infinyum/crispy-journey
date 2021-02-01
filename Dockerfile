FROM golang:1.16rc1 AS builder
WORKDIR /go/src/
COPY server.go .
RUN CGO_ENABLED=0 go build -o server server.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/server .
EXPOSE 8080
CMD ["./server"]