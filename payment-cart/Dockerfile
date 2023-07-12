FROM golang:1.20-buster as builder

# Create and change to the app directory.
WORKDIR /app

COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

ARG SKAFFOLD_GO_GCFLAGS
# Build the binary.
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath  -v -o server ./cmd

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV GOTRACEBACK=single
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Run the web service on container startup.
CMD ["app/server"]