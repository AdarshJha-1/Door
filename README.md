# Door

Door is a simple Layer 7 reverse proxy and load balancer written in Go.

I built this project to better understand how reverse proxies work internally instead of relying on libraries like `httputil.ReverseProxy`. The goal was to implement the core ideas myself and learn how requests are forwarded, how load balancing works, and how backend health is managed.

## Features

- HTTP reverse proxy
- Round Robin load balancing
- Active backend health checks
- Automatic failover
- Retry on backend failure
- Concurrent health monitoring
- Built using only Go's standard library

## Project Structure

```
cmd/
    proxy/          # Application entry point

internal/
    backend/        # Proxy and load balancing logic
    node/           # Backend node and health checks

examples/
    server/         # Example backend servers
```

## How it Works

```
            Client
               │
               ▼
          Door Proxy
               │
     ┌─────────┼─────────┐
     ▼         ▼         ▼
 Backend 1  Backend 2  Backend 3
```

1. The client sends a request to Door.
2. Door selects a healthy backend using Round Robin.
3. The request is forwarded to that backend.
4. The backend response is returned to the client.
5. If a backend fails, Door retries another healthy backend.
6. A background health checker continuously monitors backend availability.

## Running

Start the example backend servers.

```bash
go run examples/server/testserver.go
```

Start the proxy.

```bash
go run cmd/proxy/main.go
```

The proxy listens on:

```
http://localhost:6969
```

## Why I Built This

This project was built as a learning exercise to understand the internals of reverse proxies and load balancers by implementing the core functionality from scratch using Go.