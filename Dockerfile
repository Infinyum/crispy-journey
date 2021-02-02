FROM golang:1.16rc1 AS builder

RUN adduser --disabled-password robin
RUN mkdir /golang && chown robin /golang
USER robin

WORKDIR /golang
COPY . .

RUN CGO_ENABLED=0 go build -o main .

# Lightweight image for runtime only
FROM alpine:latest

RUN adduser --disabled-password robin
USER robin

COPY --from=builder /golang/main .

EXPOSE 8080
CMD ["./main"]