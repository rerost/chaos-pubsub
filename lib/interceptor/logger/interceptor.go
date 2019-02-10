package logger

import (
	"context"
	"fmt"

	"github.com/srvc/fail"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println(info.FullMethod)
		fmt.Println(req)
		result, err := handler(ctx, req)
		return result, fail.Wrap(err)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fmt.Println(info.FullMethod)
		err := handler(srv, ss)
		return fail.Wrap(err)
	}
}
