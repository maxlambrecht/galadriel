package ca

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/HewlettPackard/galadriel/pkg/common/cryptoutil"
	"github.com/HewlettPackard/galadriel/test/certtest"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmhodges/clock"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	expectedKeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageDigitalSignature
)

var (
	expectedExtendedKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
)

func TestNewCA(t *testing.T) {
	clk := clock.NewFake()
	caCert, caKey, err := certtest.CreateTestCACertificate(clk)
	require.NoError(t, err)

	// success
	config := &Config{
		RootCert: caCert,
		RootKey:  caKey,
		Clock:    clk,
	}

	ca, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, ca)
	assert.Equal(t, config.Clock, ca.clock)
}

func TestSignX509Certificate(t *testing.T) {
	clk := clock.NewFake()
	caCert, caKey, err := certtest.CreateTestCACertificate(clk)
	require.NoError(t, err)

	config := &Config{
		RootCert: caCert,
		RootKey:  caKey,
		Clock:    clk,
	}

	ca, _ := New(config)

	key, err := cryptoutil.CreateRSAKey()
	require.NoError(t, err)
	publicKey := key.Public()

	oneMinute := 1 * time.Minute

	params := X509CertificateParams{
		PublicKey: publicKey,
		TTL:       oneMinute,
		Subject: pkix.Name{
			Organization: []string{"test-org"},
			CommonName:   "test-name",
		},
	}

	cert, err := ca.SignX509Certificate(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, cert)

	assert.NotNil(t, cert.SerialNumber)
	assert.Equal(t, []string{"test-org"}, cert.Subject.Organization)
	assert.Equal(t, "test-name", cert.Subject.CommonName)
	assert.Contains(t, cert.DNSNames, "test-name")
	assert.Equal(t, publicKey, cert.PublicKey)
	assert.False(t, cert.IsCA)
	assert.True(t, cert.BasicConstraintsValid)
	assert.Equal(t, config.Clock.Now().Add(NotBeforeTolerance), cert.NotBefore)
	assert.Equal(t, config.Clock.Now().Add(oneMinute), cert.NotAfter)
	assert.Equal(t, cert.KeyUsage, expectedKeyUsage)
	assert.Equal(t, cert.ExtKeyUsage, expectedExtendedKeyUsage)
}

func TestSignJWT(t *testing.T) {
	clk := clock.NewFake()
	caCert, caKey, err := certtest.CreateTestCACertificate(clk)
	require.NoError(t, err)

	config := &Config{
		RootCert: caCert,
		RootKey:  caKey,
		Clock:    clk,
	}

	oneMinute := 1 * time.Minute

	ca, _ := New(config)

	params := JWTParams{
		Issuer:   "test-issuer",
		Subject:  spiffeid.RequireTrustDomainFromString("test-domain"),
		Audience: []string{"test-audience-1", "test-audience-2"},
		TTL:      oneMinute,
	}

	token, err := ca.SignJWT(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, token)

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) { return ca.PublicKey, nil }, jwt.WithoutClaimsValidation())

	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.NoError(t, err)
	assert.Equal(t, claims.Issuer, "test-issuer")
	assert.Equal(t, claims.Subject, "test-domain")
	assert.Contains(t, claims.Audience, "test-audience-1")
	assert.Contains(t, claims.Audience, "test-audience-2")
	assert.Equal(t, claims.IssuedAt.Time.Unix(), clk.Now().Unix())
	assert.Equal(t, claims.ExpiresAt.Time.Unix(), clk.Now().Add(oneMinute).Unix())
}
