# Multi-stage build for ultra-small container footprint (<30MB)
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/iqbal-portfolio ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Copy compiled standalone binary
COPY --from=builder /app/iqbal-portfolio /app/iqbal-portfolio

# Create persistent data directory for SQLite
RUN mkdir -p /app/data

EXPOSE 8080
ENV PORT=8080
ENV DB_PATH=/app/data/portfolio.db

VOLUME ["/app/data"]

CMD ["/app/iqbal-portfolio"]
