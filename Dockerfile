FROM golang:1.24.2 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY internal ./internal
COPY pkg ./pkg
COPY main.go Makefile ./

RUN make build

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /workspace/bin/demistio .

CMD ["./demistio"]
