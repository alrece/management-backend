package auth

import (
	"context"
	"testing"
	"time"

	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlacklist_AddAndCheck(t *testing.T) {
	rdb := testutil.NewTestRedis(t)
	blacklist := NewBlacklist(rdb)
	ctx := context.Background()

	ok, err := blacklist.IsBlacklisted(ctx, "test-jti-001")
	require.NoError(t, err)
	assert.False(t, ok)

	err = blacklist.AddJTI(ctx, "test-jti-001", 5*time.Minute)
	require.NoError(t, err)

	ok, err = blacklist.IsBlacklisted(ctx, "test-jti-001")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestBlacklist_UserRevoke(t *testing.T) {
	rdb := testutil.NewTestRedis(t)
	blacklist := NewBlacklist(rdb)
	ctx := context.Background()

	userID := int64(5001)
	issuedAt := time.Now()

	revoked, err := blacklist.IsUserRevoked(ctx, userID, issuedAt)
	require.NoError(t, err)
	assert.False(t, revoked)

	err = blacklist.RevokeByUser(ctx, userID)
	require.NoError(t, err)

	revoked, err = blacklist.IsUserRevoked(ctx, userID, issuedAt)
	require.NoError(t, err)
	assert.True(t, revoked)

	futureIssued := time.Now().Add(1 * time.Hour)
	revoked, err = blacklist.IsUserRevoked(ctx, userID, futureIssued)
	require.NoError(t, err)
	assert.False(t, revoked)
}
