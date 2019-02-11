package fault

import (
	"context"

	"google.golang.org/grpc"
)

var unaryServerInterceptorMap map[string]grpc.UnaryServerInterceptor
var streamServerInterceptorMap map[string]grpc.StreamServerInterceptor

func init() {
	unaryServerInterceptorMap = map[string]grpc.UnaryServerInterceptor{
		"/google.pubsub.v1.Publisher/Publish": Publish,
	}
	streamServerInterceptorMap = map[string]grpc.StreamServerInterceptor{
		"/google.pubsub.v1.Subscriber/StreamingPull": StreamingPull,
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if interceptor, ok := unaryServerInterceptorMap[info.FullMethod]; ok {
			return interceptor(ctx, req, info, handler)
		} else {
			return handler(ctx, req)
		}
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if interceptor, ok := streamServerInterceptorMap[info.FullMethod]; ok {
			return interceptor(srv, ss, info, handler)
		} else {
			return handler(srv, ss)
		}
	}
}

func Publish(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	// Duplicate call
	handler(ctx, req)
	return res, err
}

type streamingPullWrapper struct {
	grpc.ServerStream
}

func (s *streamingPullWrapper) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	// Duplicate message
	s.ServerStream.SendMsg(m)
	return err
}

func StreamingPull(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	wrapped := &streamingPullWrapper{ServerStream: ss}
	return handler(srv, wrapped)
}
