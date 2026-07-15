# -- Builder --
FROM golang:1.25.1 AS builder

WORKDIR /govpn

COPY go.mod ./
COPY . .

RUN go build -o bin/server ./cmd/server

# -- Runtime --
FROM alpine:3.22

WORKDIR /govpn

COPY --from=builder /govpn/bin/server .

EXPOSE 51820/udp

ENTRYPOINT ["./server"]
