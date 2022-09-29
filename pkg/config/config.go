package config

type Config struct {
	GrpcServerPort string
	RestServerPort string

	CertAuthority     string
	ServerCertificate string
	ServerKey         string
	ClientCertificate string
	ClientKey         string

	DbUrl     string
	DbPort    string
	DbName    string
	CachePort string
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
	cfg.DbUrl = dbUrl
	cfg.DbPort = dbPort
	cfg.DbName = dbName
	cfg.CachePort = cachePort
	return cfg
}
