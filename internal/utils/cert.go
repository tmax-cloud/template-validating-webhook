package utils

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"time"

	certResources "knative.dev/pkg/webhook/certificates/resources"
)

const (
	// APIServiceName is a name of APIService object
	serviceName = "template-webhook"
	CertDir     = "/tmp/cert"
)

// Create and Store certificates for webhook server
// server key / server cert is stored as file in certDir
// CA bundle is stored in ValidatingWebhookConfigurations
func CreateCert(ctx context.Context) error {
	// Make directory recursively
	if err := os.MkdirAll(CertDir, os.ModePerm); err != nil {
		return err
	}

	// Get service name and namespace
	svc := serviceName
	ns, err := Namespace()
	if err != nil {
		return err
	}

	// Create certs
	tlsKey, tlsCrt, caCrt, err := certResources.CreateCerts(ctx, svc, ns, time.Now().AddDate(1, 0, 0))
	if err != nil {
		return err
	}

	// Write certs to file
	keyPath := path.Join(CertDir, "tls.key")
	err = ioutil.WriteFile(keyPath, tlsKey, 0644)
	if err != nil {
		return err
	}

	crtPath := path.Join(CertDir, "tls.crt")
	err = ioutil.WriteFile(crtPath, tlsCrt, 0644)
	if err != nil {
		return err
	}

	caPath := path.Join(CertDir, "ca.crt")
	err = ioutil.WriteFile(caPath, caCrt, 0644)
	if err != nil {
		return err
	}

	return nil
}
