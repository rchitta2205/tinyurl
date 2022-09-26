package config

type Config struct {
	GrpcServerPort string
	RestServerPort string

	CertAuthority     string
	ServerCertificate string
	ServerKey         string
	ClientCertificate string
	ClientKey         string

	DbAddr    string
	DbName    string
	CacheAddr string
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
	cfg.DbAddr = dbAddr
	cfg.DbName = dbName
	cfg.CacheAddr = cacheAddr
	return cfg
}
