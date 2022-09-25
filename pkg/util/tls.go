package util

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
)

func LoadTLSCredentials(certAuthority, certificate, key string) (tls.Certificate, *x509.CertPool, error) {
	// Load certificate of the CA
	pemCA, err := ioutil.ReadFile(certAuthority)
	if err != nil {
		return tls.Certificate{}, nil, errors.WithStack(err)
	}

	// Certification pool to append CA's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return tls.Certificate{}, nil, errors.New("failed to add client CA's certificate")
	}

	// Load certificate and private key
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		return tls.Certificate{}, nil, errors.WithStack(err)
	}

	return cert, certPool, nil
}
