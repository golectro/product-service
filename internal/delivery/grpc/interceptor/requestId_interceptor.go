package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestIDHeader string = "x-request-id"

func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var requestID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(RequestIDHeader); len(values) > 0 && values[0] != "" {
				requestID = values[0]
			}
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		grpc.SetHeader(ctx, metadata.Pairs(RequestIDHeader, requestID))

		return handler(ctx, req)
	}
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
