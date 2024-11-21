package lib

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
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
