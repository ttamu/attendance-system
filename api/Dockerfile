FROM golang:1.23.6 AS base
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS development
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air"]

FROM base AS builder
RUN go build -o main .

FROM debian:bookworm-slim AS production
WORKDIR /root/
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/app/main .
EXPOSE 8080
CMD ["./main"]