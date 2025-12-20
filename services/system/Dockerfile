FROM ubuntu:24.04

# Install Go 1.25.4 + dev tools
RUN apt-get update && apt-get install -y \
    curl wget git ca-certificates gnupg netcat-traditional iputils-ping vim \
    && curl -fsSL https://go.dev/dl/go1.25.4.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && ln -s /usr/local/go/bin/go /usr/local/bin/go \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"
ENV GOPATH="/go"
WORKDIR /app

# Copy ALL source code FIRST
COPY . .

# Layer 1: Install oapi-codegen to /go/bin
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Layer 2: Download deps (cached)
RUN go mod download

# Layer 3: Generate OpenAPI + Build (NOW PATH includes /go/bin)
RUN cd services/system && \
    go generate ./... && \
    go build -o /app/system ./cmd/system/main.go

EXPOSE 8001
CMD ["/app/system"]
