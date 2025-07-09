package tls

import (
	"crypto/tls"
	"crypto/x509"
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
		logger.Error("Failed to append certificate", "file", certFile, "err", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
