package ssl_util

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notBefore = time.Now()
	notAfter  = time.Now().AddDate(10, 0, 0)

	ca = &x509.Certificate{
		SerialNumber: big.NewInt(69420),

		Subject:     pkix.Name{CommonName: "a"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"a"},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	cert = &x509.Certificate{
		SerialNumber: big.NewInt(42069),

		Subject:     pkix.Name{CommonName: "a"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"a"},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage: x509.KeyUsageDigitalSignature,
	}
)

type SslUtil struct {
	Dir string
}

func (s *SslUtil) GenCaFiles(caName string) error {
	caKeyPath := s.Dir + "/" + caName + ".key"
	caPath := s.Dir + "/" + caName + ".crt"

	if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
		caPubKey, caPrivKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, caPubKey, caPrivKey)
		if err != nil {
			return err
		}

		caPem, err := os.OpenFile(caPath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		pem.Encode(caPem, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		})

		caPrivKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(caPrivKey)
		if err != nil {
			return err
		}
		caPrivKeyPem, err := os.OpenFile(caKeyPath, os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		pem.Encode(caPrivKeyPem, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: caPrivKeyMarshalled,
		})
	}

	return nil
}

func (s *SslUtil) GenCertFiles(caName string, name string) error {
	caKeyPath := s.Dir + "/" + caName + ".key"
	keyPath := s.Dir + "/" + name + ".key"
	certPath := s.Dir + "/" + name + ".crt"

	if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
		caPrivKeyPem, err := os.ReadFile(caKeyPath)
		if err != nil {
			return err
		}
		caPrivKeyMarshalled, _ := pem.Decode(caPrivKeyPem)
		if caPrivKeyMarshalled == nil || caPrivKeyMarshalled.Type != "PRIVATE KEY" {
			return err
		}
		caPrivKey, err := x509.ParsePKCS8PrivateKey(caPrivKeyMarshalled.Bytes)
		if err != nil {
			return err
		}

		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		caBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, pubKey, caPrivKey)
		if err != nil {
			return err
		}

		certPem, err := os.OpenFile(certPath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		pem.Encode(certPem, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		})

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			return err
		}
		privKeyPem, err := os.OpenFile(keyPath, os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		pem.Encode(privKeyPem, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privKeyMarshalled,
		})
	}

	return nil
}

func (s *SslUtil) GenCert(caName string, pubKeyPem []byte) ([]byte, error) {
	caKeyPath := s.Dir + "/" + caName + ".key"

	caPrivKeyPem, err := os.ReadFile(caKeyPath)
	if err != nil {
		return nil, err
	}
	caPrivKeyMarshalled, _ := pem.Decode(caPrivKeyPem)
	if caPrivKeyMarshalled == nil || caPrivKeyMarshalled.Type != "PRIVATE KEY" {
		return nil, status.Error(codes.Internal, "failed to parse ca priv key")
	}
	caPrivKey, err := x509.ParsePKCS8PrivateKey(caPrivKeyMarshalled.Bytes)
	if err != nil {
		return nil, err
	}

	pubKeyMarshalled, _ := pem.Decode(pubKeyPem)
	if pubKeyMarshalled == nil || pubKeyMarshalled.Type != "PUBLIC KEY" {
		return nil, status.Error(codes.InvalidArgument, "failed to decode public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyMarshalled.Bytes)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, pubKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certPem := new(bytes.Buffer)
	pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return certPem.Bytes(), nil
}
