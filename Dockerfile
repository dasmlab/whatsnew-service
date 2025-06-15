# Stage 1: Build
FROM golang:latest AS builder

WORKDIR /workspace

ARG ARCH
ARG goproxy=https://proxy.golang.org
ENV GOPROXY=${goproxy}

# Install swag for Swagger doc gen
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY main-app/go.mod main-app/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy full source
COPY main-app/ .

# Generate Swagger docs
RUN swag init --generalInfo main.go --output docs

# Build the binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -ldflags "-s -w -extldflags '-static'" \
    -o whatsnew-service

# Stage 2: Run
FROM ubuntu:latest

# Install CA certs for TLS (required to talk to GitHub)
RUN apt-get update && apt-get install -y ca-certificates curl wget jq && rm -rf /var/lib/apt/lists/*

# Non-root for K8s security policies
USER 65532

WORKDIR /app
COPY --from=builder /workspace/whatsnew-service .
EXPOSE 10020 9200

ENTRYPOINT ["/app/whatsnew-service"]

