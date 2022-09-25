package config

const (
	// ServerAddress is the default address for the gRPC server, if no override is specified in the flags
	grpcServerAddress = "localhost:8080"
	restServerAddress = "localhost:8081"
	caPath            = "cert/ca-cert.pem"
	clientCertPath    = "cert/client-alice-cert.pem"
	clientKeyPath     = "cert/client-alice-key.pem"
	serverCertPath    = "cert/server-cert.pem"
	serverKeyPath     = "cert/server-key.pem"
)
