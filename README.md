# SnapCache

An in-memory key-value cache server implemented in Go, designed for simplicity, concurrency, and performance. The purpose of this project was mainly learning about writing protocols (RESP specifically), networking and thinking about performance optimizations.

This also gave me a good chance to learn Go more deeply.

---

## Benchmark Results

### Direct Cache Usage Comparison

| Operation                     | Time (µs) | Notes                                  |
|--------------------------------|-----------|---------------------------------------|
| Disk read (1MB file)           | 636       | Raw disk read                          |
| SET (store 1 value)            | 1,178     | Includes network and server overhead   |
| GET total (10,000 requests)    | 3,173,000 | End-to-end GETs from cache             |
| GET average per request        | 317       | ~2× faster than reading from disk      |


### End-to-End Usage (Over TCP)

| Operation                     | Time (µs) | Notes                                  |
|--------------------------------|-----------|---------------------------------------|
| Disk read (1MB file)           | 636       | Raw disk read                          |
| SET (store 1 value)            | 1,178     | Includes network and server overhead   |
| GET total (10,000 requests)    | 3,173,000 | End-to-end GETs from cache             |
| GET average per request        | 317       | ~2× faster than reading from disk      |

> These results demonstrate the performance benefit of using SnapCache as a caching layer, even over TCP connections.
> Usually, databases or blob storages are accessed over the network, making SnapCache even more advantageous.

---

## Features

- **TCP server** with support for multiple concurrent clients.
- **Supported commands**:
  - `SET key value` — Store a value in the cache.
  - `GET key` — Retrieve a value from the cache.
  - `DELETE key` — Remove a key/value pair.
  - `PING` — Health check; responds with `PONG`.
- **In-memory storage** using Go maps.
- **Concurrency-safe** with mutexes to protect shared data.
- **Simple logging** of connections, commands, and errors.

---

## Optional Enhancements (Future Work)

- TTL / key expiration with lazy deletion.
- Append-only file (AOF) persistence for durability.
- Command pipelining to reduce round-trip latency.
- Pub/Sub messaging system.
- Stats command: report total commands, memory usage, connected clients.

---

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/snapcache.git
   cd snapcache

