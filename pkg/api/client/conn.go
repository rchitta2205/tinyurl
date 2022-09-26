package client

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"time"
	"tinyurl/pkg/config"
	"tinyurl/pkg/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const (
	grpcBackOffPolicyMaxDelay = 10 * time.Second
)

type clientConn struct {
	conn *grpc.ClientConn
	ctx  context.Context
	cfg  config.Config
}

func NewClientConn(ctx context.Context, cfg config.Config) *clientConn {
	return &clientConn{
		ctx: ctx,
		cfg: cfg,
	}
}

func (c *clientConn) dial() error {
	var err error
	bc := backoff.DefaultConfig
	bc.MaxDelay = grpcBackOffPolicyMaxDelay

	// Load client side tls credentials
	cert, certPool, err := util.LoadTLSCredentials(c.cfg.CertAuthority, c.cfg.ClientCertificate, c.cfg.ClientKey)
	if err != nil {
		return errors.WithStack(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: bc,
		}),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	c.conn, err = grpc.DialContext(c.ctx, c.cfg.GrpcServerPort, opts...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *clientConn) getGrpcConn() *grpc.ClientConn {
	return c.conn
}

func (c *clientConn) release() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return errors.WithStack(err)
		}
		c.conn = nil
	}
	return nil
}
