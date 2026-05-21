package auth

import (
	"context"
	"testing"

	"management-backend/internal/config"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginLock_FailAndLock(t *testing.T) {
	rdb := testutil.NewTestRedis(t)
	ctx := context.Background()

	config.C.LoginSecurity.MaxFailCount = 3
	lock := NewLoginLock(rdb)

	for i := 0; i < 2; i++ {
		err := lock.RecordFail(ctx, "testuser")
		require.NoError(t, err)
	}
	locked, err := lock.IsLocked(ctx, "testuser")
	require.NoError(t, err)
	assert.False(t, locked)

	err = lock.RecordFail(ctx, "testuser")
	require.NoError(t, err)
	locked, err = lock.IsLocked(ctx, "testuser")
	require.NoError(t, err)
	assert.True(t, locked)
}

func TestLoginLock_ResetFail(t *testing.T) {
	rdb := testutil.NewTestRedis(t)
	lock := NewLoginLock(rdb)
	ctx := context.Background()

	err := lock.RecordFail(ctx, "resetuser")
	require.NoError(t, err)

	err = lock.ResetFail(ctx, "resetuser")
	require.NoError(t, err)

	count, err := lock.GetFailCount(ctx, "resetuser")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestLoginLock_GetFailCount_Empty(t *testing.T) {
	rdb := testutil.NewTestRedis(t)
	lock := NewLoginLock(rdb)
	ctx := context.Background()

	count, err := lock.GetFailCount(ctx, "newuser")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
