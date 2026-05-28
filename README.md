# nam-tcp

## Introduction

**nam-tcp** is a TCP proxy service with token-based access control. It exposes a management HTTP API to register users, define backend targets (host and port), and issue short-lived personal access tokens (PATs). Clients connect to the proxy with a PAT; the server validates the token and forwards traffic to the configured target.

The project ships as two Go binaries: a **proxy server** (HTTP API + TCP proxy) and a **client** for connecting, sending data, or running a local port forward.

## What is it used for

Use **nam-tcp** when you need controlled, auditable TCP access to internal services without opening those services directly to a network. Typical use cases include:

- **Secure tunneling** to databases, message queues, or other TCP backends behind a bastion or jump host
- **Time-limited access** via PATs that map a user to a specific target and expire automatically
- **Local port forwarding** so tools like `psql`, `redis-cli`, or any TCP client can talk to a remote host through the proxy as if it were listening on `127.0.0.1`

Flow at a glance:

1. Register **users** and **targets** (e.g. `localhost:5432`) through the HTTP API.
2. Create a **user PAT** linking a user to a target; save the plaintext token (shown only once).
3. Run the **client** with that token; the proxy authenticates and bridges TCP to the target.

## Tech stack

| Layer | Technology |
|-------|------------|
| Language | [Go](https://go.dev/) 1.26+ |
| HTTP API | `net/http` (Go 1.22+ route patterns) |
| TCP proxy | Standard library `net` |
| Database | [SQLite](https://www.sqlite.org/) via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) |
| IDs | [google/uuid](https://github.com/google/uuid) |
| Auth tokens | SHA-256 hashed PATs (`nam_tcp_…` prefix) |

Project layout:

- `cmd/proxy` — HTTP server (`:8000`) and TCP proxy (`:8888`)
- `cmd/client` — CLI: `forward`, `connect`, `send`
- `internal/` — handlers, services, repositories, proxy protocol
- `Dockerfile.server` / `Dockerfile.client` — multi-stage Docker builds
- `Makefile` — build and run images locally

## How to use

### Prerequisites

- [Go](https://go.dev/) 1.26 or newer (for local development)
- [Docker](https://www.docker.com/) (for container builds and `make` targets)
- Clone this repository and run commands from the project root

You can run **nam-tcp** with `go run` / `go build`, or with Docker via the **Makefile** (recommended for a consistent environment).

### Makefile and Docker

List available targets:

```bash
make help
```

| Target | Description |
|--------|-------------|
| `make build` | Build `nam-tcp-server` and `nam-tcp-client` images |
| `make build-server` | Build server image only |
| `make build-client` | Build client image only |
| `make run-server` | Run proxy (HTTP `:8000`, TCP `:8888`, SQLite volume `nam-tcp-data`) |
| `make run-client` | Run client; pass subcommand via `ARGS` |
| `make clean` | Remove built images |

Optional variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_IMAGE` | `nam-tcp-server` | Server image name |
| `CLIENT_IMAGE` | `nam-tcp-client` | Client image name |
| `IMAGE_TAG` | `latest` | Image tag |
| `ARGS` | _(empty)_ | Arguments passed to `nam-tcp-client` |
| `NAM_TCP_PROXY` | — | Default proxy address for the client container |
| `NAM_TCP_TOKEN` | — | Default PAT for the client container |
| `NAM_TCP_LOCAL` | — | Default local bind for `forward` |

**Start the server in Docker:**

```bash
make run-server
```

**Run the client in Docker** (uses host networking so `localhost:8888` reaches the server on your machine):

```bash
make run-client ARGS="forward -local 127.0.0.1:15432 -proxy localhost:8888 -token nam_tcp_<your-token>"
```

Or set env vars once:

```bash
export NAM_TCP_PROXY=localhost:8888
export NAM_TCP_TOKEN=nam_tcp_<your-token>
make run-client ARGS="connect"
```

Custom image names:

```bash
make build SERVER_IMAGE=myregistry/nam-tcp-server IMAGE_TAG=v1.0.0
```

### 1. Start the proxy server

**With Make (Docker):**

```bash
make run-server
```

**With Go (local):**

```bash
go run ./cmd/proxy
```

Environment variables (optional):

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8000` | Management API listen address |
| `PROXY_ADDR` | `:8888` | TCP proxy listen address |

On startup, the server creates `app.db` (SQLite) and runs schema migrations.

### 2. Configure users, targets, and PATs

**Health check**

```bash
curl http://localhost:8000/health
```

**Create a user**

```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com", "password":"secret"}'
```

**Create a target**

```bash
curl -X POST http://localhost:8000/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name":"postgres", "host":"localhost", "port": "5432"}'
```

**Create a PAT** (use `user_id` and `target_id` from the responses above)

```bash
curl -X POST http://localhost:8000/api/v1/user-pats \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<user-uuid>", "target_id":"<target-uuid>", "ttl_in_hour": 12}'
```

The response includes `plaintext` (e.g. `nam_tcp_…`). Store it securely; only the hash is kept in the database.

**Other API endpoints**

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/users` | List users |
| `GET` | `/api/v1/targets` | List targets |
| `GET` | `/api/v1/user-pats` | List PATs |
| `DELETE` | `/api/v1/user-pats/{id}` | Revoke a PAT |
| `DELETE` | `/api/v1/user-pats/expired` | Revoke all expired PATs |

### 3. Use the client

Set defaults with environment variables if you prefer:

| Variable | Description |
|----------|-------------|
| `NAM_TCP_PROXY` | Proxy address (default `localhost:8888`) |
| `NAM_TCP_TOKEN` | PAT token |
| `NAM_TCP_LOCAL` | Local bind address for `forward` (default `127.0.0.1:15432`) |

**Forward** — listen locally and tunnel each connection through the proxy to the PAT’s target:

```bash
# Docker
make run-client ARGS="forward -local 127.0.0.1:15432 -proxy localhost:8888 -token nam_tcp_<your-token>"

# Go
go run ./cmd/client forward \
  -local 127.0.0.1:15432 \
  -proxy localhost:8888 \
  -token "nam_tcp_<your-token>"
```

Then point your tool at the local port, for example:

```bash
psql -h 127.0.0.1 -p 15432 -U admin mydb
```

**Connect** — interactive stdin/stdout session through the proxy:

```bash
make run-client ARGS="connect -proxy localhost:8888 -token nam_tcp_<your-token>"

# or: go run ./cmd/client connect -proxy localhost:8888 -token "nam_tcp_<your-token>"
```

**Send** — send a payload and print the response:

```bash
make run-client ARGS="send -proxy localhost:8888 -token nam_tcp_<your-token> -data hello"

# or: go run ./cmd/client send -proxy localhost:8888 -token "nam_tcp_<your-token>" -data "hello"
```

You can pass `-token-file /path/to/token` instead of `-token`.

### 4. Example end-to-end workflow

This walkthrough uses Docker for the server and client. API setup uses `curl` against the server on `localhost:8000`. See `scripts/play.sh` for more commented examples.

**Terminal 1 — start the server**

```bash
make run-server
```

The server listens on HTTP `:8000` and TCP proxy `:8888`. SQLite data is stored in the Docker volume `nam-tcp-data`.

**Terminal 2 — register resources and create a PAT**

```bash
# Health check
curl http://localhost:8000/health

# Create user (save id from response)
curl -s -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com", "password":"secret"}'

# Create target pointing at a TCP service reachable from the server container
# For a DB on the host machine, use host.docker.internal (Docker Desktop) or the host IP
curl -s -X POST http://localhost:8000/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name":"postgres", "host":"host.docker.internal", "port": "5432"}'

# Create PAT (replace user_id and target_id; save plaintext token from response)
curl -s -X POST http://localhost:8000/api/v1/user-pats \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<user-uuid>", "target_id":"<target-uuid>", "ttl_in_hour": 12}'
```

**Terminal 3 — forward a local port through the proxy**

```bash
export NAM_TCP_TOKEN=nam_tcp_<plaintext-from-create-response>

make run-client ARGS="forward -local 127.0.0.1:15432 -proxy localhost:8888 -token $NAM_TCP_TOKEN"
```

**Terminal 4 — use your TCP client**

```bash
psql -h 127.0.0.1 -p 15432 -U admin mydb
```

**Optional — quick connectivity test**

```bash
make run-client ARGS="send -proxy localhost:8888 -token $NAM_TCP_TOKEN -data ping"
```

#### Local development (without Docker)

```bash
# Terminal 1
go run ./cmd/proxy

# Terminal 2 — curl setup (use host "localhost" for targets on the same machine)
# ...

# Terminal 3
go run ./cmd/client forward -local 127.0.0.1:15432 -proxy localhost:8888 -token "$NAM_TCP_TOKEN"
```

### Build binaries

**Docker images:**

```bash
make build
```

**Go binaries:**

```bash
go build -o nam-tcp-proxy ./cmd/proxy
go build -o nam-tcp-client ./cmd/client
```

## License

See repository license terms if applicable.
