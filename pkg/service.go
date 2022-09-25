package pkg

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"sync"
	"tinyurl/pkg/api/proto"
	"tinyurl/pkg/api/server"
	"tinyurl/pkg/config"
	"tinyurl/pkg/interceptor"
	"tinyurl/pkg/util"
)

type Service interface {
	Register() error
	Serve(*sync.WaitGroup) (net.Listener, error)
}

type grpcService struct {
	ctx        context.Context
	grpcServer *grpc.Server
	cfg        config.Config
	logEntry   *logrus.Entry
}

type restService struct {
	ctx        context.Context
	restServer *http.Server
	cfg        config.Config
	logEntry   *logrus.Entry
}

func NewGrpcService(ctx context.Context, cfg config.Config, logEntry *logrus.Entry) Service {
	return &grpcService{
		ctx:      ctx,
		cfg:      cfg,
		logEntry: logEntry,
	}
}

func NewRestService(ctx context.Context, cfg config.Config, logEntry *logrus.Entry) Service {
	return &restService{
		ctx:      ctx,
		cfg:      cfg,
		logEntry: logEntry,
	}
}

func (s *grpcService) Register() error {
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
	authInterceptor := interceptor.NewAuthInterceptor(appMgr.AuthApplication(), s.logEntry)

	// Add all the unary interceptors
	unaryInterceptors = append(unaryInterceptors, authInterceptor.UnaryAuthInterceptor)

	// Add all the stream interceptors
	streamInterceptors = append(streamInterceptors, authInterceptor.StreamAuthInterceptor)

	// Create tls credentials
	cert, certPool, err := util.LoadTLSCredentials(s.cfg.CertAuthority, s.cfg.ClientCertificate, s.cfg.ClientKey)
	if err != nil {
		return errors.WithStack(err)
	}

	// Configure credentials to require and verify the client cert
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}

	// Create the grpc server object
	s.grpcServer = grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// Register the grpc app servers
	proto.RegisterTinyUrlServiceServer(s.grpcServer, tinyUrlServer)

	s.logEntry.Info("Registered Grpc servers")
	return nil
}

func (s *grpcService) Serve(wg *sync.WaitGroup) (net.Listener, error) {
	grpcLis, err := net.Listen("tcp", s.cfg.GrpcServerAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	reflection.Register(s.grpcServer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logEntry.Info("Starting gRPC server...")
		err = s.grpcServer.Serve(grpcLis)
		s.logEntry.Warn("Grpc Server: " + err.Error())
	}()

	return grpcLis, errors.WithStack(err)
}

func (s *restService) Register() error {
	// Creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()

	// Reverse proxy Rest service calls the gRPC client so needs the client certificates
	cert, certPool, err := util.LoadTLSCredentials(s.cfg.CertAuthority, s.cfg.ClientCertificate, s.cfg.ClientKey)
	if err != nil {
		return errors.WithStack(err)
	}

	// Configure tlsConfig
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
	}

	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.DefaultConfig,
		}),
	}

	// Setting up a dial-up for gRPC service by specifying endpoint/target url
	err = proto.RegisterTinyUrlServiceHandlerFromEndpoint(s.ctx, mux, s.cfg.GrpcServerAddress, options)
	if err != nil {
		return errors.WithStack(err)
	}

	// Creating an HTTP server
	s.restServer = &http.Server{
		Handler: mux,
	}

	s.logEntry.Info("Registered http servers")
	return nil
}

func (s *restService) Serve(wg *sync.WaitGroup) (net.Listener, error) {
	httpLis, err := net.Listen("tcp", s.cfg.RestServerAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logEntry.Info("Starting REST server...")
		err = s.restServer.Serve(httpLis)
		s.logEntry.Warn("Rest Server: " + err.Error())
	}()

	return httpLis, errors.WithStack(err)
}
