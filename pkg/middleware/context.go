package middleware

import "context"

type contextKey string

const (
	keyUserID    contextKey = "user_id"
	keyUsername  contextKey = "username"
	keyTenantID  contextKey = "tenant_id"
	keyRequestID contextKey = "request_id"
)

func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

func GetUserID(ctx context.Context) int64 {
	if v, ok := ctx.Value(keyUserID).(int64); ok {
		return v
	}
	return 0
}

func WithUsername(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, keyUsername, name)
}

func GetUsername(ctx context.Context) string {
	if v, ok := ctx.Value(keyUsername).(string); ok {
		return v
	}
	return ""
}

func WithTenantID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, keyTenantID, id)
}

func GetTenantID(ctx context.Context) int64 {
	if v, ok := ctx.Value(keyTenantID).(int64); ok {
		return v
	}
	return 0
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(keyRequestID).(string); ok {
		return v
	}
	return ""
}
