# Functional Requirements for ProtoCache

## Core Functional Requirements

### 1. TCP Server
- The system must run as a standalone TCP server.
- Must accept multiple client connections concurrently.
- Must parse incoming protobuf-encoded commands.

### 2. Supported Commands

#### **SET**
- Inputs: key, value, optional TTL.
- Stores the entry in memory (overwrites existing).
- Applies expiration if TTL is provided.

#### **GET**
- Input: key.
- Returns the stored value if present and not expired.
- Returns `NOT_FOUND` if missing or expired.

#### **DELETE**
- Input: key.
- Removes the entry from the cache.
- Succeeds even if the key does not exist (idempotent).

#### **PING**
- Returns `PONG` for health checks.

### 3. In-Memory Storage
- Uses a Go map or sync.Map as the key-value store.
- Each entry stores:
  - key  
  - value  
  - expiration timestamp (nullable)

### 4. Expiration Handling
- Expired keys must not be returned.
- Lazy deletion:
  - Expiration is checked during `GET`.
  - If expired, delete it and return `NOT_FOUND`.
- Optional periodic cleanup goroutine.

### 5. Protobuf Serialization
- All network communication uses protobuf messages.
- Must define `.proto` files for:
  - Request (command type, key, value, ttl)
  - Response (status, value)

### 6. Concurrency Safety
- Server must handle concurrent operations safely.
- Use mutex or RWMutex around access to the store.

### 7. Logging
- Log:
  - New client connections
  - Commands received
  - Errors
  - Key expiration events (optional)

### 8. Error Handling
- Handle:
  - Invalid protobuf input
  - Unknown commands
  - Empty keys
  - Connection failures

---

## Optional Enhancements

### AOF Persistence
- Log each SET/DELETE to a file.
- Replay on startup.

### Command Pipelining
- Support batches of commands in a single message.

### Pub/Sub
- Allow clients to subscribe and receive messages.

### Stats Command
- Expose:
  - total commands processed
  - memory usage approximation
  - number of connected clients

