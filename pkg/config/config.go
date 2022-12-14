package config

type Config struct {
	GrpcServerPort string
	RestServerPort string

	CertAuthority     string
	ServerCertificate string
	ServerKey         string
	ClientCertificate string
	ClientKey         string

	Db    string
	Cache string

	DaprAddr string
	DaprPort string
}

func NewConfig() Config {
	cfg := Config{}
	cfg.GrpcServerPort = grpcServerPort
	cfg.RestServerPort = restServerPort
	cfg.CertAuthority = caPath
	cfg.ServerCertificate = serverCertPath
	cfg.ServerKey = serverKeyPath
	cfg.ClientCertificate = clientCertPath
	cfg.ClientKey = clientKeyPath
	cfg.Db = db
	cfg.Cache = cache
	cfg.DaprAddr = daprAddr
	cfg.DaprPort = daprPort
	return cfg
}
