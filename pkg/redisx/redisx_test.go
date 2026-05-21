package redisx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) (*goredis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close(); mr.Close() })
	return client, mr
}

// --- Lock Tests ---

func TestLockAcquireAndUnlock(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()

	lock, err := AcquireLock(ctx, client, "test-key", DefaultLockOpts())
	if err != nil {
		t.Fatalf("acquire lock: %v", err)
	}

	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("unlock: %v", err)
	}
}

func TestLockMutualExclusion(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()

	lock1, err := AcquireLock(ctx, client, "mutex", LockOpts{TTL: 5 * time.Second, RetryCount: 1, RetryDelay: 10 * time.Millisecond})
	if err != nil {
		t.Fatalf("acquire lock1: %v", err)
	}

	_, err = AcquireLock(ctx, client, "mutex", LockOpts{TTL: 5 * time.Second, RetryCount: 2, RetryDelay: 10 * time.Millisecond})
	if err == nil {
		t.Fatal("should not acquire second lock")
	}

	lock1.Unlock(ctx)
}

func TestLockExpiryOnCancel(t *testing.T) {
	client, mr := setupTestRedis(t)
	ctx := context.Background()

	opts := LockOpts{TTL: 500 * time.Millisecond, RetryCount: 1, RetryDelay: 10 * time.Millisecond}
	lock, err := AcquireLock(ctx, client, "expiry", opts)
	if err != nil {
		t.Fatalf("acquire lock: %v", err)
	}

	lock.cancel()
	mr.FastForward(600 * time.Millisecond)

	lock2, err := AcquireLock(ctx, client, "expiry", opts)
	if err != nil {
		t.Fatalf("should acquire lock after expiry: %v", err)
	}
	lock2.Unlock(ctx)
}

func TestLockNonOwnerUnlock(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()

	lock, _ := AcquireLock(ctx, client, "owner-test", LockOpts{TTL: 5 * time.Second, RetryCount: 1, RetryDelay: 10 * time.Millisecond})
	lock.value = "wrong-value"

	err := lock.Unlock(ctx)
	if err == nil {
		t.Fatal("non-owner should fail to unlock")
	}
}

// --- Cache Tests ---

func TestCacheSetAndGet(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	cache := NewCache(client, "test:")

	type User struct{ Name string }
	err := Set(ctx, cache, "user:1", User{Name: "alice"}, 10*time.Minute)
	if err != nil {
		t.Fatalf("cache set: %v", err)
	}

	result, err := Get[User](ctx, cache, "user:1")
	if err != nil {
		t.Fatalf("cache get: %v", err)
	}
	if result.Name != "alice" {
		t.Fatalf("expected alice, got %s", result.Name)
	}
}

func TestCacheMiss(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	cache := NewCache(client, "test:")

	_, err := Get[string](ctx, cache, "nonexistent")
	if !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}
}

func TestCacheDelete(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	cache := NewCache(client, "test:")

	Set(ctx, cache, "del-key", "value", 10*time.Minute)
	Delete(ctx, cache, "del-key")

	_, err := Get[string](ctx, cache, "del-key")
	if !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("expected ErrCacheMiss after delete, got %v", err)
	}
}

// --- Idempotent Tests ---

func TestIdempotentFirstExecution(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	exec := NewIdempotentExecutor(client, "idem:")

	called := 0
	result, err := Execute(ctx, exec, "op-1", func() (string, error) {
		called++
		return "done", nil
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result != "done" {
		t.Fatalf("expected done, got %s", result)
	}
	if called != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestIdempotentDuplicateExecution(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	exec := NewIdempotentExecutor(client, "idem:")

	called := 0
	fn := func() (string, error) {
		called++
		return "done", nil
	}

	Execute(ctx, exec, "op-2", fn)
	Execute(ctx, exec, "op-2", fn)

	if called != 1 {
		t.Fatalf("expected 1 call (idempotent), got %d", called)
	}
}

func TestIdempotentFailureAllowsRetry(t *testing.T) {
	client, _ := setupTestRedis(t)
	ctx := context.Background()
	exec := NewIdempotentExecutor(client, "idem:")

	called := 0
	fn := func() (string, error) {
		called++
		if called == 1 {
			return "", errors.New("transient error")
		}
		return "recovered", nil
	}

	_, err := Execute(ctx, exec, "op-3", fn)
	if err == nil {
		t.Fatal("first call should fail")
	}

	result, err := Execute(ctx, exec, "op-3", fn)
	if err != nil {
		t.Fatalf("retry should succeed: %v", err)
	}
	if result != "recovered" {
		t.Fatalf("expected recovered, got %s", result)
	}
	if called != 2 {
		t.Fatalf("expected 2 calls, got %d", called)
	}
}
