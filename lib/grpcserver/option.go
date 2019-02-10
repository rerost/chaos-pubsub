package grpcserver

import (
	"google.golang.org/grpc"
)

type Option func(c Config) Config

type Config struct {
	Servers                  []Server
	Port                     int64
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor
	StreamServerInterceptors []grpc.StreamServerInterceptor
}

func WithServers(servers ...Server) Option {
	return func(c Config) Config {
		c.Servers = append(c.Servers, servers...)
		return c
	}
}

func WithPort(port int64) Option {
	return func(c Config) Config {
		c.Port = port
		return c
	}
}

func WithGrpcServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(c Config) Config {
		c.UnaryServerInterceptors = append(c.UnaryServerInterceptors, interceptors...)
		return c
	}
}

func WithGrpcServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(c Config) Config {
		c.StreamServerInterceptors = append(c.StreamServerInterceptors, interceptors...)
		return c
	}
}
