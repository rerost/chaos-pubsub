package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/srvc/fail"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := map[string]interface{}{}
		log["type"] = "unary"
		log["time"] = time.Now().Unix()
		log["method"] = info.FullMethod
		log["request"] = req
		j, err := json.Marshal(log)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))

		result, err := handler(ctx, req)
		return result, fail.Wrap(err)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := map[string]interface{}{}
		log["type"] = "stream"
		log["time"] = time.Now().Unix()
		log["method"] = info.FullMethod

		j, err := json.Marshal(log)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))

		err = handler(srv, ss)
		return fail.Wrap(err)
	}
}
