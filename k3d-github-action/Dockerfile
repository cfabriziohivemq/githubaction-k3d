FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "Building on $BUILDPLATFORM for $TARGETPLATFORM"

WORKDIR /app
COPY . .

# Set the correct GOARCH based on target platform
RUN case "$TARGETPLATFORM" in \
    "linux/amd64") export GOARCH=amd64 ;; \
    "linux/arm64") export GOARCH=arm64 ;; \
    *) export GOARCH=amd64 ;; \
    esac && \
    go build -o /bin/k3d-action ./cmd/action

FROM alpine:3.18

# Install required tools
RUN apk add --no-cache docker-cli curl

# Install k3d
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | sh

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && mv kubectl /usr/local/bin/

COPY --from=builder /bin/k3d-action /usr/local/bin/k3d-action

ENTRYPOINT ["k3d-action"]
