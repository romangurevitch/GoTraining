# Stage 1: Build
# Use the official Go image as the builder. Pin to a specific version for reproducibility.
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching.
# The module download step is only re-run when go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary targeting Linux amd64.
# -trimpath:       removes local file system paths from the binary (reproducible builds)
# -ldflags="-w -s": strips debug info and symbol table, reducing binary size
# CGO_ENABLED=0:   disables cgo so the binary is fully static (no libc dependency)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-w -s" \
    -o /hello ./cmd/hello/main.go


# Stage 2: Runtime
# Use a minimal scratch image — contains nothing but our binary.
# This keeps the final image as small as possible and reduces the attack surface.
FROM scratch

# Copy only the compiled binary from the builder stage
COPY --from=builder /hello /hello

# Document the default entry point
ENTRYPOINT ["/hello"]