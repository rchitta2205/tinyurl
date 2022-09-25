package config

type Config struct {
	GrpcServerAddress string
	RestServerAddress string

	CertAuthority     string
	ServerCertificate string
	ServerKey         string

	ClientCertificate string
	ClientKey         string
}

func NewConfig() Config {
	cfg := Config{}
	cfg.GrpcServerAddress = grpcServerAddress
	cfg.RestServerAddress = restServerAddress
	cfg.CertAuthority = caPath
	cfg.ServerCertificate = serverCertPath
	cfg.ServerKey = serverKeyPath
	cfg.ClientCertificate = clientCertPath
	cfg.ClientKey = clientKeyPath
	return cfg
}
