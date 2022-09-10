/*
 * The package implements wrappers for grpc server configurations
 */

package tlsconf

/*
 * TLS Configurations Package
 */

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

/*
 * Public Functions
 */

// CertPoolWithCustomRootCA appends a user defined custom CA to the system's CAs and returns the new certificate pool
func CertPoolWithCustomRootCA(customCertFile string, logger *logrus.Logger) (*x509.CertPool, error) {

	var certPool *x509.CertPool
	var err error
	// load the Root CAs on the system
	if certPool, err = x509.SystemCertPool(); err != nil {
		return nil, err
	}

	// parse the the certificate file
	// do not fail hard in case the user doesn't specifiy the file
	var pemFileBytes []byte
	if pemFileBytes, err = ioutil.ReadFile(customCertFile); err != nil {
		logger.Warnf("Failed to load the certificate file: %s", customCertFile)
		return certPool, nil
	}

	// load cert from the PEM file
	if ok := certPool.AppendCertsFromPEM(pemFileBytes); !ok {
		logger.Warnf("Invalid certificate file not added to certificate pool: %s", customCertFile)
		return certPool, nil

	}

	// all loaded CAs
	return certPool, nil
}

// ClientTLSConfigWithCustomRootCA loads a custom CA, appends it to Root CAs, and then return a client TLS configuration
func ClientTLSConfigWithCustomRootCA(customCertFile string, logger *logrus.Logger) (*tls.Config, error) {

	// get the cert pool
	certPool, err := CertPoolWithCustomRootCA(customCertFile, logger)
	if err != nil {
		return nil, err
	}

	// create config
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: false, // verify remote peer/server

	}
	return tlsConfig, nil
}

/*
 * Private Functions
 */
