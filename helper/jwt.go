package helper

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type JWTClaim struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GetJWTKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GenerateToken(userId int64) (string, error) {
	claims := &JWTClaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTKey())
}

func VerifyToken(tokenStr string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GRPCWithAuth(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func ForwardMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(context.Background(), md)
}
