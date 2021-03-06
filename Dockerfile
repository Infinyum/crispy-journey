FROM golang:1.16 AS builder

RUN mkdir /golang

WORKDIR /golang
COPY src/crispy/ .

RUN CGO_ENABLED=0 go build -o main .

# Lightweight image for runtime only
FROM alpine:latest

COPY --from=builder /golang/static static/
COPY --from=builder /golang/main .

EXPOSE 8080
CMD ["./main"]