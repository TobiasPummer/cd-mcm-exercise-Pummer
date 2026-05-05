# Docker Documentation – Product Catalog API

## Dockerfile Analysis

### Multi-Stage Build

The Dockerfile uses a two-stage build:

**Stage 1 – builder:**
```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api
```
This stage uses the full Go toolchain to compile the binary.
Dependencies are downloaded before copying the source code — this
way Docker can cache the dependency layer and only re-downloads
when go.mod or go.sum change.

**Stage 2 – runtime:**
```dockerfile
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /api-server .
EXPOSE 8080
ENTRYPOINT ["./api-server"]
```
Only the compiled binary is copied from the builder stage into a
minimal Alpine image. The Go toolchain, source code, and
intermediate build files are discarded.

### CGO_ENABLED=0

`CGO_ENABLED=0` disables C bindings in the Go compiler. This produces
a fully statically linked binary with no external library dependencies.

This is important because:
- The binary can run in minimal containers like Alpine or even scratch
- No C runtime libraries need to be present in the final image
- The image is smaller and has a reduced attack surface

### Image Size Comparison

| Build type | Approximate size |
|---|---|
| Single-stage (golang:alpine) | ~300 MB |
| Multi-stage (alpine runtime) | ~20 MB |

The multi-stage build reduces the image size by roughly 93% by
excluding the Go toolchain from the final image.

---

## Docker Compose

The `docker-compose.yml` orchestrates two services:

- **db**: PostgreSQL 16 with a named volume `pgdata` for persistence
- **api**: The Go API built from the Dockerfile, depends on db being healthy

The API only starts after the database passes its healthcheck
(`pg_isready`), preventing connection errors on startup.

---

## CRUD Test Results

### Setup
```bash
docker compose up --build
```

### Health Check
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### Create Products
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":999.99}'
# {"id":1,"name":"Laptop","price":999.99}

curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Maus","price":29.99}'
# {"id":2,"name":"Maus","price":29.99}

curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Tastatur","price":59.99}'
# {"id":3,"name":"Tastatur","price":59.99}
```

### List All Products
```bash
curl http://localhost:8080/products
# [{"id":1,"name":"Laptop","price":999.99},{"id":2,"name":"Maus","price":29.99},{"id":3,"name":"Tastatur","price":59.99}]
```

### Update Product
```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Pro","price":1299.99}'
# {"id":1,"name":"Laptop Pro","price":1299.99}
```

### Delete Product
```bash
curl -X DELETE http://localhost:8080/products/3
# {"result":"success"}
```

### List After Delete
```bash
curl http://localhost:8080/products
# [{"id":1,"name":"Laptop Pro","price":1299.99},{"id":2,"name":"Maus","price":29.99}]
```

---

## Data Persistence Test

After stopping and restarting the containers:

```bash
docker compose down && docker compose up
curl http://localhost:8080/products
# [{"id":1,"name":"Laptop Pro","price":1299.99},{"id":2,"name":"Maus","price":29.99}]
```

All products are still present after restart. The named volume
`pgdata` in docker-compose.yml ensures PostgreSQL data persists
across container restarts.