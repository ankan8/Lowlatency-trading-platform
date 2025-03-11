package middleware

import (
    "context"
    "errors"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v4"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

func UnaryJWTInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {

    
    fullMethod := info.FullMethod

    //Skip token check for public endpoints

    if strings.HasSuffix(fullMethod, "Register") || strings.HasSuffix(fullMethod, "Login") {
        
        return handler(ctx, req)
    }

    //Otherwise, do the normal token validation
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, errors.New("missing metadata")
    }
    authHeaders := md["authorization"]
    if len(authHeaders) == 0 {
        return nil, errors.New("missing authorization token")
    }
    tokenStr := authHeaders[0]

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, errors.New("JWT_SECRET is not set")
    }
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    
    return handler(ctx, req)
}
