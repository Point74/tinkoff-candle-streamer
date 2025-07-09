package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
)

func LoadTLSCredentials(certFile string, logger *slog.Logger) (credentials.TransportCredentials, error) {
	cert, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		logger.Error("Failed to append certificate", "file", certFile)
		return nil, fmt.Errorf("failed to append certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
