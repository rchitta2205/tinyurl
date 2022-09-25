package pkg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"tinyurl/pkg/api/proto"
	"tinyurl/pkg/api/server"
	"tinyurl/pkg/auth"
	"tinyurl/pkg/config"
)

type service struct {
	ctx      context.Context
	server   *grpc.Server
	cfg      config.Config
	logEntry *logrus.Entry
}

func NewService(ctx context.Context, cfg config.Config, logEntry *logrus.Entry) *service {
	return &service{
		ctx:      ctx,
		cfg:      cfg,
		logEntry: logEntry,
	}
}

func (s *service) Register() error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	// Create the application manager
	appMgr, err := NewApplicationManagerBuilder(s.logEntry).Build()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Create all the app servers
	tinyUrlServer, err := server.NewTinyUrlServer(appMgr.TinyUrlApplication())
	if err != nil {
		return errors.WithStack(err)
	}

	// Create all the interceptors
	authInterceptor := auth.NewAuthInterceptor(appMgr.AuthApplication(), s.logEntry)

	// Add all the unary interceptors
	unaryInterceptors = append(unaryInterceptors, authInterceptor.UnaryAuthInterceptor)

	// Add all the stream interceptors
	streamInterceptors = append(streamInterceptors, authInterceptor.StreamAuthInterceptor)

	// Create tls credentials
	cred, err := s.loadTLSCredentials()
	if err != nil {
		return errors.WithStack(err)
	}

	// Create the grpc server object
	s.server = grpc.NewServer(
		grpc.Creds(cred),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// Register the app servers
	proto.RegisterTinyUrlServiceServer(s.server, tinyUrlServer)
	return nil
}

func (s *service) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(s.cfg.CertAuthority)
	if err != nil {
		return nil, err
	}

	// Certification pool to append client CA's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(s.cfg.ServerCertificate, s.cfg.ServerKey)
	if err != nil {
		return nil, err
	}

	// Configure credentials to require and verify the client cert
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func (s *service) Serve() error {
	lis, err := net.Listen("tcp", s.cfg.ServerAddress)
	if err != nil {
		return errors.WithStack(err)
	}
	defer lis.Close()

	reflection.Register(s.server)

	err = s.server.Serve(lis)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
