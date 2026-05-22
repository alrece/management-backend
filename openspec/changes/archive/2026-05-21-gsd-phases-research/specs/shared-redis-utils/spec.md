## ADDED Requirements

### Requirement: Distributed Lock
The system SHALL provide a Redis-based distributed lock using SetNX with UUID value, background TTL extension, and Lua atomic unlock.

#### Scenario: Acquire lock successfully
- **WHEN** caller invokes Lock(ctx, key, ttl) and key is not held
- **THEN** system sets key with UUID value and TTL, starts background renewal goroutine, returns Lock struct

#### Scenario: Acquire lock on already-locked key
- **WHEN** caller invokes Lock(ctx, key, ttl) and key is already held by another
- **THEN** system retries with exponential backoff until context timeout, returns error if not acquired

#### Scenario: Unlock by lock owner
- **WHEN** caller invokes Unlock() on a held lock and UUID matches
- **THEN** system deletes key via Lua script (check UUID then DEL), stops renewal goroutine

#### Scenario: Unlock by non-owner
- **WHEN** caller invokes Unlock() and UUID does not match
- **THEN** system returns error without deleting key

#### Scenario: Lock auto-expiry on holder crash
- **WHEN** lock holder crashes and renewal stops
- **THEN** key expires after TTL, other instances can acquire the lock

### Requirement: Cache Wrapper
The system SHALL provide typed cache operations with TTL support using cache-aside pattern.

#### Scenario: Cache hit
- **WHEN** caller invokes Get(ctx, key) and key exists in Redis
- **THEN** system deserializes value into typed result, returns nil error

#### Scenario: Cache miss
- **WHEN** caller invokes Get(ctx, key) and key does not exist
- **THEN** system returns ErrCacheMiss error

#### Scenario: Set with TTL
- **WHEN** caller invokes Set(ctx, key, value, ttl)
- **THEN** system serializes value and stores in Redis with specified TTL

#### Scenario: Delete cached key
- **WHEN** caller invokes Delete(ctx, key)
- **THEN** system removes key from Redis regardless of existence

### Requirement: Idempotent Execution
The system SHALL provide idempotent execution guard using Redis SetNX with 24h TTL.

#### Scenario: First execution
- **WHEN** caller invokes Execute(ctx, key, fn) for a new key
- **THEN** system acquires idempotent key, executes fn, stores result, returns result

#### Scenario: Duplicate execution within TTL
- **WHEN** caller invokes Execute(ctx, key, fn) for an existing key
- **THEN** system returns cached result without executing fn

#### Scenario: TTL expiry
- **WHEN** idempotent key expires after 24h
- **THEN** subsequent call treats as first execution
