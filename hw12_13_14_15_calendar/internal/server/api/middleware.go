package internalapi

import (
	"context"
	"fmt"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func loggingMiddleware(logg logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		t1 := time.Now()
		resp, err := handler(ctx, req)
		timing := time.Since(t1)
		var remoteAddr, userAgent string
		if meta, ok := metadata.FromIncomingContext(ctx); ok {
			headers := meta.Get("user-agent")
			if len(headers) > 0 {
				userAgent = headers[0]
			}
		}
		if p, ok := peer.FromContext(ctx); ok {
			remoteAddr = p.Addr.String()
		}
		if err != nil {
			logg.Error(
				fmt.Sprintf(
					`%s [%s] %s %.2f "%s" "%s"`,
					remoteAddr,
					time.Now().Format(time.RFC3339),
					info.FullMethod,
					float32(timing/time.Second),
					err.Error(),
					userAgent,
				),
			)
		} else {
			logg.Info(
				fmt.Sprintf(
					`%s [%s] %s %.2f "%s"`,
					remoteAddr,
					time.Now().Format(time.RFC3339),
					info.FullMethod,
					float32(timing/time.Second),
					userAgent,
				),
			)
		}
		return resp, err
	}
}
