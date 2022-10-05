package config

const (
	// ServerAddress is the default address for the gRPC server, if no override is specified in the flags
	grpcServerPort = ":8081"
	restServerPort = ":8080"
	caPath         = "cert/ca-cert.pem"
	clientCertPath = "cert/client-alice-cert.pem"
	clientKeyPath  = "cert/client-alice-key.pem"
	serverCertPath = "cert/server-cert.pem"
	serverKeyPath  = "cert/server-key.pem"
	db             = "mongo"
	cache          = "redis"
	daprAddr       = "127.0.0.1"
	daprPort       = "50001"
)
