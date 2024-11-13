package ssl_util

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os"
)

var (
	ca = &x509.Certificate{
		SerialNumber: big.NewInt(69420),

		Subject:     pkix.Name{CommonName: "a"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"a"},

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	cert = &x509.Certificate{
		SerialNumber: big.NewInt(42069),

		Subject:     pkix.Name{CommonName: "a"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"a"},

		KeyUsage: x509.KeyUsageDigitalSignature,
	}
)

type SslUtil struct {
	dir string
}

func (s *SslUtil) GenCaFiles(caName string) {
	caKeyPath := s.dir + "/" + caName + ".key"
	caPath := s.dir + "/" + caName + ".crt"

	if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
		caPubKey, caPrivKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatalf("failed to generate ca keypair: %v", err)
		}
		caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, caPubKey, caPrivKey)
		if err != nil {
			log.Fatalf("failed to create ca cert: %v", err)
		}

		caPem, err := os.Create(caPath)
		if err != nil {
			log.Fatalf("failed to create %s: %v", caPath, err)
		}
		pem.Encode(caPem, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		})

		caPrivKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(caPrivKey)
		if err != nil {
			log.Fatalf("failed to marshal ca priv key: %v", err)
		}
		caPrivKeyPem, err := os.Create(caKeyPath)
		if err != nil {
			log.Fatalf("failed to create %s: %v", caKeyPath, err)
		}
		pem.Encode(caPrivKeyPem, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: caPrivKeyMarshalled,
		})
	}
}

func (s *SslUtil) GenCertFiles(caName string, name string) {
	caKeyPath := s.dir + "/" + caName + ".key"
	keyPath := s.dir + "/" + name + ".key"
	certPath := s.dir + "/" + name + ".crt"

	if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
		caPrivKeyPem, err := os.ReadFile(caKeyPath)
		if err != nil {
			log.Fatalf("failed to read ca.key: %v", err)
		}
		caPrivKeyMarshalled, _ := pem.Decode(caPrivKeyPem)
		if caPrivKeyMarshalled == nil || caPrivKeyMarshalled.Type != "PRIVATE KEY" {
			log.Fatal("failed to decode ca.key")
		}
		caPrivKey, err := x509.ParsePKCS8PrivateKey(caPrivKeyMarshalled.Bytes)
		if err != nil {
			log.Fatalf("failed to parse ca.key: %v", err)
		}

		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatalf("failed to generate keypair: %v", err)
		}
		bytes, err := x509.CreateCertificate(rand.Reader, cert, ca, pubKey, caPrivKey)
		if err != nil {
			log.Fatalf("failed to create cert: %v", err)
		}

		certPem, err := os.Create(certPath)
		if err != nil {
			log.Fatalf("failed to create server.crt: %v", err)
		}
		pem.Encode(certPem, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: bytes,
		})

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			log.Fatalf("failed to marshal priv key: %v", err)
		}
		privKeyPem, err := os.Create(keyPath)
		if err != nil {
			log.Fatalf("failed to create server.key: %v", err)
		}
		pem.Encode(privKeyPem, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privKeyMarshalled,
		})
	}
}
