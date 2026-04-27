# Go Toolchain & Docker

Modern Go development revolves around its powerful, built-in toolchain and containerisation for production-grade deployments.

---

## 1. Core Concepts

| Tool / File | Description / Purpose |
| :--- | :--- |
| **Go Modules** | Go's dependency management system. Replaces the legacy GOPATH workspace. |
| **`go.mod`** | The manifest file for your module, defining its name and requirements. |
| **`go.sum`** | A lockfile containing checksums of specific module versions. |
| **`Dockerfile`** | A script for building container images. Multi-stage builds are the Go standard. |

---

## 2. Visual Representation

```mermaid
flowchart TD
    S["Source Code\n(.go files)"] --"go build"--> B["Binary\n(Statically Linked)"]
    B --"COPY"--> D["Docker Runtime Image\n(Minimal)"]
    D --> R["Production Container\n(Small & Secure)"]
```

---

## 3. Project Structure

This module is a **standalone Go project** with its own `go.mod`. It demonstrates a simple CLI app with configuration management.

```
toolchain/
├── main.go           # Entry point with greeting logic and Viper config (commented out)
├── main_test.go      # Unit tests for the greeting function
├── go.mod            # Module definition with Viper dependency
├── go.sum            # Dependency checksums
├── Dockerfile        # Multi-stage Docker build
├── config/
│   ├── dev.yaml      # Development environment config
│   └── prod.yaml     # Production environment config
└── README.md
```

---

## 4. Implementation Examples

### Essential CLI Commands

Initialise a new module:
```bash
go mod init my-cool-project # Run this when starting a new project
```

Add missing dependencies and remove unused ones:
```bash
go mod tidy
```

Run your application directly (from the toolchain directory):
```bash
go run . Gopher
```

Compile a binary for the current OS/Arch:
```bash
go build -o hello-demo .
```

Run tests:
```bash
go test -v ./...
```

Install a tool to your `$GOPATH/bin`:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Adding a Dependency

The `main.go` file contains commented-out code that uses [Viper](https://github.com/spf13/viper) for configuration management.
To enable it:

1. Install the dependency:
```bash
go get github.com/spf13/viper
```

2. Uncomment the Viper import and the `loadConfig()` function in `main.go`

3. Run the application to see the config values loaded from `config/dev.yaml`

```bash
go run .
# Output: App: hello-dev, Port: 8080, LogLevel: debug
#         Hello, Go Bank!
```

This demonstrates how Go manages external dependencies — adding an import triggers `go mod tidy` to fetch and verify the package.

### Multi-Stage Dockerfile
```dockerfile
# Stage 1: Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /app/myapp .

# Stage 2: Final Runtime
FROM scratch
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/config /config
ENTRYPOINT ["/myapp"]
```

The second `COPY` brings the `config/` directory into the runtime image so the application can load its YAML configuration files at runtime. Without it, the `scratch` image would contain only the binary and config loading would fail.

---

## 5. Common Patterns & Use Cases

- **Caching Dependencies**: Copying `go.mod` and `go.sum` before the rest of the source in a Dockerfile to cache the `go mod download` layer.
- **Reproducible Builds**: Using `-trimpath` and pinned Go versions in Docker to ensure the same source always produces the same binary.
- **Minimal Images**: Using `scratch` or `alpine` for the final runtime image to reduce the attack surface and image size.

---

## 6. Critical Pitfalls & Best Practices

> **Warning:** Never commit your vendor directory or local configuration files to source control. Rely on `go.mod` and environment variables.

1. **Keep go.mod Tidy**: Frequently run `go mod tidy` to prune unused dependencies.
2. **Checksum Integrity**: Never manually edit `go.sum`. It is managed automatically by the toolchain to ensure supply-chain security.
3. **Static Binaries**: Always set `CGO_ENABLED=0` when building for `scratch` or `distroless` images to avoid runtime failures due to missing C libraries.

---

## Running the Examples

Navigate to the toolchain directory and explore the Go toolchain:

```bash
cd internal/basics/toolchain
```

1. Build the example:
```bash
go build -o hello-demo .
```

2. Run the compiled binary:
```bash
./hello-demo
./hello-demo Gopher
```

3. Run the tests:
```bash
go test -v ./...
```

4. Build the Docker image:
```bash
docker build -t hello:latest .
```

5. Run the container:
```bash
docker run --rm hello:latest Gopher
```

## Your Next Step
Now that you can build and package your code, it's time to dive into the core language primitives.
Explore **[Basic Types](../types/README.md)** to learn about variables, slices, and maps.

---

## Further Reading

- [Go Modules Reference](https://go.dev/ref/mod)
- [Viper Configuration Library](https://github.com/spf13/viper)
- [Dockerizing a Go Application](https://docs.docker.com/language/golang/build-images/)
- [Shrinking your Go binaries](https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/)
