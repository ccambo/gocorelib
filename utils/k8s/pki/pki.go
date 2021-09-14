/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package pki

import (
	"crypto"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
)

// ===== [ Constants and Variables ] =====
const (
	rsaKeySize = 2048
)

var ()

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewPrivateKey - RSA 알고리즘을 사용하는 개인키 생성
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
}

// NewCSR - 지정한 정보를 기준으로 CSR 생성
func NewCSR(config certutil.Config, key crypto.Signer) (*x509.CertificateRequest, error) {
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   config.CommonName,
			Organization: config.Organization,
		},
		DNSNames:    config.AltNames.DNSNames,
		IPAddresses: config.AltNames.IPs,
	}

	csrBytes, err := x509.CreateCertificateRequest(cryptorand.Reader, template, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a CSR")
	}

	return x509.ParseCertificateRequest(csrBytes)
}

// NewCSRAndKey - 지정한 정보를 기준으로 인증서의 서명을 생성할 수 있는 새로운 키와 CSR (Certificate Signing Request - 인증서 서명 요청) 생성
func NewCSRAndKey(config *certutil.Config) (*x509.CertificateRequest, *rsa.PrivateKey, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create private key")
	}

	csr, err := NewCSR(*config, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to generate CSR")
	}

	return csr, key, nil
}
