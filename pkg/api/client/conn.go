package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
	"tinyurl/pkg/config"

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
	tlsCredentials, err := c.loadTLSCredentials()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: bc,
		}),
		grpc.WithTransportCredentials(tlsCredentials),
	}

	c.conn, err = grpc.DialContext(c.ctx, c.cfg.ServerAddress, opts...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *clientConn) getGrpcConn() *grpc.ClientConn {
	return c.conn
}

func (c *clientConn) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(c.cfg.CertAuthority)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate %v", pemServerCA)
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(c.cfg.ClientCertificate, c.cfg.ClientKey)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
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
