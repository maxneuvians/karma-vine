FROM mcr.microsoft.com/devcontainers/go:1.25

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    nodejs \
    npm && \
    rm -rf /var/lib/apt/lists/*

# Install air (live-reload)
RUN go install github.com/air-verse/air@latest

RUN go install golang.org/x/tools/gopls@latest && go install golang.org/x/tools/cmd/goimports@latest && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

