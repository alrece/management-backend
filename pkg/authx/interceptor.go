package authx

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKeyClaims struct{}

// GRPCAuthInterceptor gRPC 认证拦截器（从 metadata 提取 token → 解析验证 → 注入 context）
func GRPCAuthInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "缺少 metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "缺少 authorization")
		}

		tokenStr := tokens[0]
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := ParseToken(tokenStr, jwtSecret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Token 验证失败: %v", err)
		}

		ctx = context.WithValue(ctx, ctxKeyClaims{}, claims)
		return handler(ctx, req)
	}
}

// ClaimsFromContext 从 gRPC context 提取 Claims
func ClaimsFromContext(ctx context.Context) *Claims {
	if c, ok := ctx.Value(ctxKeyClaims{}).(*Claims); ok {
		return c
	}
	return nil
}
