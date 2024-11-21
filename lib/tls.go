package lib

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
	notAfter  = time.Now().AddDate(1000, 0, 0)

	ca = &x509.Certificate{
		SerialNumber: big.NewInt(69420),

		Subject:     pkix.Name{CommonName: "nf6"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"nf6"},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	cert = &x509.Certificate{
		SerialNumber: big.NewInt(42069),

		Subject:     pkix.Name{CommonName: "nf6"},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"nf6"},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage: x509.KeyUsageDigitalSignature,
	}
)

func GenCaFiles(dir string, caName string) error {
	caKeyPath := dir + "/" + caName + ".key"
	caPath := dir + "/" + caName + ".crt"

	if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
		caPubKey, caPrivKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return status.Error(codes.Internal, "failed to generate ca private key")
		}
		caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, caPubKey, caPrivKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to create ca certificate")
		}

		caPrivKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(caPrivKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to marshall ca private key")
		}
		caPrivKeyPem, err := os.OpenFile(caKeyPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return status.Error(codes.Internal, "failed to open ca private key file")
		}
		if err := pem.Encode(caPrivKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: caPrivKeyMarshalled}); err != nil {
			return status.Error(codes.Internal, "failed to encode ca private key")
		}

		caPem, err := os.OpenFile(caPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return status.Error(codes.Internal, "failed to open ca certificate file")
		}
		if err := pem.Encode(caPem, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
			return status.Error(codes.Internal, "failed to encode ca certificate")
		}
	}

	return nil
}

func GenCertFiles(dir string, caName string, name string) error {
	caKeyPath := dir + "/" + caName + ".key"
	keyPath := dir + "/" + name + ".key"
	certPath := dir + "/" + name + ".crt"

	if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
		caPrivKeyPem, err := os.ReadFile(caKeyPath)
		if err != nil {
			return status.Error(codes.Internal, "failed to read ca certificate file")
		}
		caPrivKeyMarshalled, _ := pem.Decode(caPrivKeyPem)
		if caPrivKeyMarshalled == nil || caPrivKeyMarshalled.Type != "PRIVATE KEY" {
			return status.Error(codes.Internal, "failed to decode ca private key")
		}
		caPrivKey, err := x509.ParsePKCS8PrivateKey(caPrivKeyMarshalled.Bytes)
		if err != nil {
			return status.Error(codes.Internal, "failed to parse ca private key")
		}

		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return status.Error(codes.Internal, "failed to generate keypair")
		}

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to marshall private key")
		}
		privKeyPem, err := os.OpenFile(keyPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return status.Error(codes.Internal, "failed to open private key file")
		}
		if err := pem.Encode(privKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyMarshalled}); err != nil {
			return status.Error(codes.Internal, "failed to encode private key")
		}

		caBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, pubKey, caPrivKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to create certificate")
		}
		certPem, err := os.OpenFile(certPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return status.Error(codes.Internal, "failed to open certificate file")
		}
		if err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
			return status.Error(codes.Internal, "failed to encode certificate")
		}
	}

	return nil
}

func GenCert(dir string, caName string, pubKeyPem []byte) ([]byte, error) {
	caKeyPath := dir + "/" + caName + ".key"

	caPrivKeyPem, err := os.ReadFile(caKeyPath)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to read ca private key file")
	}
	caPrivKeyMarshalled, _ := pem.Decode(caPrivKeyPem)
	if caPrivKeyMarshalled == nil || caPrivKeyMarshalled.Type != "PRIVATE KEY" {
		return nil, status.Error(codes.Internal, "failed to decode ca private key")
	}
	caPrivKey, err := x509.ParsePKCS8PrivateKey(caPrivKeyMarshalled.Bytes)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse ca private key")
	}

	pubKeyMarshalled, _ := pem.Decode(pubKeyPem)
	if pubKeyMarshalled == nil || pubKeyMarshalled.Type != "PUBLIC KEY" {
		return nil, status.Error(codes.InvalidArgument, "failed to decode public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyMarshalled.Bytes)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse public key")
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, pubKey, caPrivKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create certificate")
	}

	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return nil, status.Error(codes.Internal, "failed to encode certificate")
	}

	return certPem.Bytes(), nil
}

func GenKeyFiles(dir string, name string) error {
	privKeyPath := dir + "/" + name + ".key"
	pubKeyPath := dir + "/" + name + ".pub"

	if _, err := os.Stat(pubKeyPath); errors.Is(err, os.ErrNotExist) {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return status.Error(codes.Internal, "failed to generate keypair")
		}

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to marshall private key")
		}
		privKeyPem, err := os.OpenFile(privKeyPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return status.Error(codes.Internal, "failed to open private key file")
		}
		if err := pem.Encode(privKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyMarshalled}); err != nil {
			return status.Error(codes.Internal, "failed to encode private key")
		}

		pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			return status.Error(codes.Internal, "failed to marshall public key")
		}
		pubKeyPem, err := os.OpenFile(pubKeyPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return status.Error(codes.Internal, "failed to open public key file")
		}
		if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
			return status.Error(codes.Internal, "failed to encode public key")
		}
	} else {
		return status.Error(codes.AlreadyExists, "keypair already exists")
	}

	return nil
}
