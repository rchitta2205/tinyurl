package config

type Config struct {
	ServerAddress string

	CertAuthority     string
	ServerCertificate string
	ServerKey         string

	ClientCertificate string
	ClientKey         string
}

func NewConfig() Config {
	cfg := Config{}
	cfg.ServerAddress = serverAddress
	cfg.CertAuthority = caPath
	cfg.ServerCertificate = serverCertPath
	cfg.ServerKey = serverKeyPath
	cfg.ClientCertificate = clientCertPath
	cfg.ClientKey = clientKeyPath
	return cfg
}

