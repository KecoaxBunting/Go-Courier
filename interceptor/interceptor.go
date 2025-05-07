package auth

import (
	"context"
	"fmt"
	"strings"

	helper "go-courier/helper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ContextKey string

const UserIdKey ContextKey = "user"

func AuthInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Login" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
			return nil, fmt.Errorf("missing or invalid authorization")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		claims, err := helper.VerifyToken(token)
		if err != nil {
			return nil, fmt.Errorf("invalid or expired token")
		}

		newCtx := context.WithValue(ctx, UserIdKey, claims.UserId)
		return handler(newCtx, req)
	}
}
