# Builder stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /workspace

# copy go.mod/go.sum first to take advantage of layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/system ./services/system/cmd/system/main.go

# Runtime: Ubuntu with full features
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y \
    curl wget netcat-traditional iputils-ping \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app

# final image: small, non-root runtime
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/system /system

EXPOSE 8001
USER nonroot

ENTRYPOINT ["/system"]
# ...existing code...