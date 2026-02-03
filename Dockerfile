# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api-server main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /app

COPY --from=builder /app/api-server /app/api-server
RUN chmod +x /app/api-server && ls -la /app/

EXPOSE 8080

ENTRYPOINT ["/app/api-server"]
