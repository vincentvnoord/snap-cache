# SnapCache

An in-memory key-value cache server implemented in Go, designed for simplicity, concurrency, and performance. The purpose of this project was mainly learning about writing protocols (RESP specifically), networking and thinking about performance optimizations.

This also gave me a good chance to learn Go more deeply.

---

## Benchmark Results

### Direct Cache Usage Comparison
When using the SnapCache package directly in Go (no TCP, no serialization, no syscalls), operations are extremely fast and limited only by CPU and Go runtime overhead.

| Operation                     | Time | Notes                                  |
|--------------------------------|-----------|---------------------------------------|
| Disk read (1MB file, 1 time)   | ~450 µs      | Raw disk read                          |
| SET average                    | ~70 ns      | Direct memory read                  |
| GET total (sum of 10,000 calls)       | ~300 µs    | GETs from direct cache             |
| GET average                    | ~30 ns      | ~15,000× faster than reading from disk      |

> These values represent pure in-memory access, not networked cache behavior.
> This highlights the theoretical maximum speed of the cache layer.

### End-to-End Usage (Over TCP)

| Operation                     | Time (µs) | Notes                                  |
|--------------------------------|-----------|---------------------------------------|
| Disk read (1MB file, 1 time)   | ~450      | Raw disk read                          |
| SET (store 1 value)            | ~1,500     | Network + server + serialization  |
| GET total (sum of 10,000 requests)    | ~4,000,000 | End-to-end GETs over network             |
| GET average per request        | ~400       | ~1.12× faster than reading from disk      |

> These results demonstrate the performance benefit of using SnapCache as a caching layer, even over TCP connections.
> In this test the server directly reads from its local file system, however, usually databases or blob storages are accessed over the network, making SnapCache even more advantageous.

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

