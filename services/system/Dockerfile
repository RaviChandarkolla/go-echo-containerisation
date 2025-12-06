FROM golang:1.25.4 AS builder

WORKDIR /workspace

# copy go.mod/go.sum first to take advantage of layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/system ./services/system/cmd/system/main.go

# final image: small, non-root runtime
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/system /system

EXPOSE 8001
USER nonroot

ENTRYPOINT ["/system"]
# ...existing code...