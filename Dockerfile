FROM golang:1.16rc1 AS builder
# Secure against running as root
RUN adduser --disabled-password robin
RUN mkdir /golang && chown robin /golang
USER robin
WORKDIR /golang
COPY server.go .
RUN CGO_ENABLED=0 go build -o server server.go

FROM alpine:latest
# Secure against running as root
RUN adduser --disabled-password robin
USER robin
COPY --from=builder /golang/server .
EXPOSE 8080
CMD ["./server"]