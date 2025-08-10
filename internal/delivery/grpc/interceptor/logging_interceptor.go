package interceptor

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		method := info.FullMethod

		requestID := GetRequestID(ctx)

		var clientIP string
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		resp, err := handler(ctx, req)

		latency := time.Since(start)

		statusCode := "OK"
		level := logrus.InfoLevel

		if err != nil {
			st, _ := status.FromError(err)
			statusCode = st.Code().String()

			switch st.Code() {
			case codes.Internal, codes.DataLoss, codes.Unavailable:
				level = logrus.ErrorLevel
			case codes.InvalidArgument, codes.Unauthenticated, codes.PermissionDenied:
				level = logrus.WarnLevel
			default:
				level = logrus.InfoLevel
			}
		}

		fields := logrus.Fields{
			"method":    method,
			"status":    statusCode,
			"latency":   latency,
			"requestId": requestID,
			"clientIP":  clientIP,
		}

		entry := logger.WithFields(fields)

		msg := "gRPC request handled"
		entry.Log(level, msg)

		return resp, err
	}
}
