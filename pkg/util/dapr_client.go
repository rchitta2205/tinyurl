package util

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func Connect(ctx context.Context, addr string, opts ...grpc.DialOption) (dapr.Client, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 60 * time.Second,
			Backoff:           backoff.DefaultConfig,
		})}

	dialOpts = append(dialOpts, opts...)
	var client dapr.Client
	grpcConn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err == nil {
		client = dapr.NewClientWithConnection(grpcConn)
		return client, nil
	}

	return nil, errors.Wrapf(err, "Failed to open Dapr connection: %v", err)
}
