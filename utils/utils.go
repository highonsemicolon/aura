package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func InitTLS(CACertPath string) (*tls.Config, error) {
	caCert, err := os.ReadFile(CACertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}

	return &tls.Config{RootCAs: rootCertPool}, nil
}
