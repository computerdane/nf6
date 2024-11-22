package lib

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	TlsName = "nf6"

	notBefore = time.Now()
	notAfter  = time.Now().AddDate(1000, 0, 0)

	TlsCaTemplate = &x509.Certificate{
		SerialNumber: big.NewInt(69420),

		Subject:     pkix.Name{CommonName: TlsName},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{TlsName},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	TlsCertTemplate = &x509.Certificate{
		SerialNumber: big.NewInt(42069),

		Subject:     pkix.Name{CommonName: TlsName},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{TlsName},
		NotBefore:   notBefore,
		NotAfter:    notAfter,

		KeyUsage: x509.KeyUsageDigitalSignature,
	}
)

func TlsGenKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to generate key")
	}
	return pubKey, privKey, nil
}

func TlsEncodePrivKey(privKey ed25519.PrivateKey) ([]byte, error) {
	privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshall private key")
	}
	privKeyPem := new(bytes.Buffer)
	if err := pem.Encode(privKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyMarshalled}); err != nil {
		return nil, status.Error(codes.Internal, "failed to encode private key")
	}
	return privKeyPem.Bytes(), nil
}

func TlsEncodePubKey(pubKey ed25519.PublicKey) ([]byte, error) {
	pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshall public key")
	}
	pubKeyPem := new(bytes.Buffer)
	if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
		return nil, status.Error(codes.Internal, "failed to encode public key")
	}
	return pubKeyPem.Bytes(), nil
}

func TlsDecodePrivKey(privKeyPem []byte) (ed25519.PrivateKey, error) {
	privKeyMarshalled, _ := pem.Decode(privKeyPem)
	if privKeyMarshalled == nil || privKeyMarshalled.Type != "PRIVATE KEY" {
		return nil, status.Error(codes.Internal, "failed to decode private key")
	}
	privKeyParsed, err := x509.ParsePKCS8PrivateKey(privKeyMarshalled.Bytes)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse private key")
	}
	privKey, ok := privKeyParsed.(ed25519.PrivateKey)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to cast private key")
	}
	return privKey, nil
}

func TlsDecodePubKey(pubKeyPem []byte) (ed25519.PublicKey, error) {
	pubKeyMarshalled, _ := pem.Decode(pubKeyPem)
	if pubKeyMarshalled == nil || pubKeyMarshalled.Type != "PUBLIC KEY" {
		return nil, status.Error(codes.Internal, "failed to decode public key")
	}
	pubKeyParsed, err := x509.ParsePKIXPublicKey(pubKeyMarshalled.Bytes)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse public key")
	}
	pubKey, ok := pubKeyParsed.(ed25519.PublicKey)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to cast public key")
	}
	return pubKey, nil
}

func TlsWriteFile(data []byte, path string) error {
	if _, err := os.Stat(path); err == nil {
		return status.Error(codes.AlreadyExists, "tls file already exists")
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return status.Error(codes.Internal, "failed to open tls file")
	}
	file.Write(data)
	return nil
}

func TlsReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to read tls file")
	}
	return data, nil
}

func TlsGenCert(template *x509.Certificate, pubKey ed25519.PublicKey, caPrivKey ed25519.PrivateKey) ([]byte, error) {
	data, err := x509.CreateCertificate(rand.Reader, template, TlsCaTemplate, pubKey, caPrivKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create certificate")
	}
	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: data}); err != nil {
		return nil, status.Error(codes.Internal, "failed to encode certificate")
	}
	return certPem.Bytes(), nil
}

func TlsGenCertUsingPrivKeyFile(template *x509.Certificate, pubKey ed25519.PublicKey, caPrivKeyPath string) ([]byte, error) {
	caPrivKeyPem, err := TlsReadFile(caPrivKeyPath)
	if err != nil {
		return nil, err
	}
	caPrivKey, err := TlsDecodePrivKey(caPrivKeyPem)
	if err != nil {
		return nil, err
	}
	return TlsGenCert(template, pubKey, caPrivKey)
}
