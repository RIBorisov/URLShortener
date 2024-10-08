package interceptors

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"shortener/internal/models"
	"shortener/internal/service"
)

// UserIdUnaryInterceptor generates and adds JWT token from metadata to context.
func UserIdUnaryInterceptor(svc *service.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			token, err := svc.BuildJWTString(svc.SecretKey)
			if err != nil {
				svc.Log.Err("Failed to generate token", err)
				return nil, status.Error(codes.Unauthenticated, "Access denied")
			}
			md = metadata.New(map[string]string{"token": token})
		}
		token := md.Get("token")
		if len(token) == 0 {
			generatedToken, err := svc.BuildJWTString(svc.SecretKey)
			if err != nil {
				svc.Log.Err("Failed to generate token", err)
				return nil, status.Error(codes.Unauthenticated, "Access denied")
			}
			md.Set("token", generatedToken)
			token = append(token, generatedToken)
		}

		userID := svc.GetUserID(token[0], svc.SecretKey, svc.Log)
		if userID == "" {
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}
		newCtx := context.WithValue(ctx, models.CtxUserIDKey, userID)

		return handler(newCtx, req)
	}
}

var ErrTokenNotExists = errors.New("failed to get token from metadata")
