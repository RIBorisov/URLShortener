package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"shortener/internal/logger"
)

// main generates new cert and key.
func main() {
	log := &logger.Log{}
	log.Initialize("DEBUG")
	log.Debug("run generating..")
	if err := prepareTLS(log); err != nil {
		log.Fatal("failed to prepare TLS cert and key", err)
	}
	log.Debug("success!")
}

func prepareTLS(log *logger.Log) error {
	const (
		certPath string = "tls/server.crt"
		keyPath  string = "tls/server.key"
	)
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create cert: %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = certFile.Close()
		if err != nil {
			log.Fatal("failed to close file", err)
		}
	}()
	if _, err = certFile.Write(certPEM.Bytes()); err != nil {
		return fmt.Errorf("failed to write cert file: %w", err)
	}

	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = keyFile.Close()
		if err != nil {
			log.Fatal("failed to close file", err)
		}
	}()
	if _, err = keyFile.Write(privateKeyPEM.Bytes()); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}
