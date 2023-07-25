package grpcpool

import (
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"time"
)

func DefaultGrpcCallOpts() []grpc.CallOption {
	return []grpc.CallOption{
		grpc_retry.WithMax(1),
		grpc_retry.WithPerRetryTimeout(time.Second * 5),
	}
}

var DefaultOptions = Options{
	Init:            32,
	MaxActive:       256,
	IdleTimeout:     time.Minute * 3,
	MaxLifeDuration: time.Minute * 3,
}
