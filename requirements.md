# SnapCache

A in-memory key-value cache server implemented in Go.

## Features

- **TCP server** with concurrent client support.
- **Commands**:
  - `SET key value` — store a value.
  - `GET key` — retrieve a value.
  - `DELETE key` — remove a value.
  - `PING` — health check (`PONG` response).
- **In-memory storage** using Go maps.
- **Concurrency-safe** with mutexes.
- **Simple logging** of connections, commands, and errors.

## Optional Enhancements

- TTL / key expiration with lazy deletion.
- Append-only file (AOF) persistence.
- Command pipelining for batch requests.
- Pub/Sub messaging.
- Stats command: total commands, memory usage, connected clients.
