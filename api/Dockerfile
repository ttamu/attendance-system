FROM golang:1.23.6 AS base
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS development
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM base AS builder

WORKDIR /go/src/app/cmd/server
RUN go build -o /go/bin/server

WORKDIR /go/src/app/cmd/seed
RUN go build -o /go/bin/seed

FROM debian:bookworm-slim AS production
WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/server .
COPY --from=builder /go/bin/seed .

EXPOSE 8080

CMD ["sh", "-c", "./server & ./seed & wait"]