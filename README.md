# SnapCache

An in-memory key-value cache server implemented in Go, designed for simplicity, concurrency, and performance. The purpose of this project was mainly learning about writing protocols (RESP specifically), networking and thinking about performance optimizations.

This also gave me a good chance to learn Go more deeply.

---

## Benchmark Results


### End-to-End Usage Comparison (Over TCP)
The following benchmarks compare database operations with and without using SnapCache as a caching layer over TCP. The case involves a user and their orders, where we frequently need to fetch the sum of order amounts per user. The test (in /cmd/benchmark/database/main.go) simulates fetching this data 100,000 times, measuring the time taken for each operation.

| Operation                     | Min Time     | Max Time      | Avg Time     | Total Time       | Notes                                |
|--------------------------------|-------------|---------------|-------------|-----------------|--------------------------------------|
| DB: SELECT Username + ID       | 171.228µs   | 1.645545ms    | 219.593µs   | 2.195938s       | Simple select by ID                  |
| DB: SELECT User order SUM      | 55.799ms    | 70.709ms      | 58.892ms    | 5.88925s        | Aggregation with join                |
| CACHE: SET user order sum      | 19.26µs     | 2.985873ms    | 24.819µs    | 1.984978s       | Storing aggregated sums in cache     |
| CACHE: GET user order sum      | 16.31µs     | 2.028561ms    | 24.886µs    | 1.990396s       | Fetching aggregated sums from cache  |
| Fetched from DB / Cache        | 79,978      | —             | —           | —               | Number of records retrieved          |

> Considering the average times, the cache layer is approximately 2,368× faster than fetching the aggregated sums directly from the database. This shows the significant performance improvement that caching can provide for read-heavy operations.

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

