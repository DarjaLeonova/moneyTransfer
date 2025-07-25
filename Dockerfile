# === Build Stage ===
FROM golang:1.23.0-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Force Go to build a Linux binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.* ./
RUN go mod tidy

COPY . ./

RUN go build -o server cmd/main.go

# === Run Stage ===
FROM alpine:latest

WORKDIR /app

# Copy compiled Linux binary
COPY --from=builder /app/server .

# Copy environment file (correctly!)
COPY .env /app



# Run the binary
CMD ["sh", "-c", "env && ./server"]
