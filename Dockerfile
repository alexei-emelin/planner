FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o app main.go

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/web ./web
COPY --from=builder /app/app .
COPY --from=builder /app/.env .env

RUN apt-get update && apt-get install -y --no-install-recommends \
    libsqlite3-0 ca-certificates && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 7540

CMD ["./app"]
