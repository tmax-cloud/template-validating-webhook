package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

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

func ReadCertFile(certFile string) ([]byte, error) {
	bytedCertFile, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	return bytedCertFile, nil
}

func UpdateCABundle(caCrt string) {
	webhook := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	s := scheme.Scheme

	if err := admissionregistrationv1.AddToScheme(s); err != nil {
		panic(err)
	}
	c, err := Client(client.Options{Scheme: s})
	if err != nil {
		panic(err)
	}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: "template-validate-webhook", Namespace: ""}, webhook); err != nil {
		panic(err)
	}

	updateWebhook := webhook.DeepCopy()
	bytedCert, _ := ReadCertFile(caCrt)

	// when updating CABundle field, it is automatically encoded as base64
	if webhook != nil {
		updateWebhook.Webhooks[0].ClientConfig.CABundle = bytedCert
	}

	if err := c.Patch(context.TODO(), updateWebhook, client.MergeFrom(webhook)); err != nil {
		fmt.Println(err)
	}
}
