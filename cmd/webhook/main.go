package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tmax-cloud/template-validating-webhook/internal/utils"
	"github.com/tmax-cloud/template-validating-webhook/pkg/apis"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	cert   = utils.CertDir + "/tls.crt"
	key    = utils.CertDir + "/tls.key"
	caCert = utils.CertDir + "/ca.crt"
)

var log = logf.Log.WithName("template-validating-webhook")

func main() {
	log.Info("initializing server....")

	if err := utils.CreateCert(context.Background()); err != nil {
		fmt.Println(err, "failed to create cert")
	}
	//TODO whenever pod restart, update CABundle of ValidatingWebhookConfiguration CRD
	webhook := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	s := scheme.Scheme

	if err := admissionregistrationv1.AddToScheme(s); err != nil {
		panic(err)
	}

	c, err := utils.Client(client.Options{Scheme: s})
	if err != nil {
		panic(err)
	}

	if err := c.Get(context.TODO(), types.NamespacedName{Name: "template-validate-webhook", Namespace: ""}, webhook); err != nil {
		panic(err)
	}

	updateWebhook := webhook.DeepCopy()
	bytedCert, _ := utils.ReadCertFile(caCert)

	fmt.Println(string(bytedCert))

	if webhook != nil {
		updateWebhook.Webhooks[0].ClientConfig.CABundle = bytedCert
	}

	if err := c.Patch(context.TODO(), updateWebhook, client.MergeFrom(webhook)); err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/validate", apis.CheckInstanceUpdatable).Methods("POST")

	http.Handle("/", r)

	if err := http.ListenAndServeTLS(":8443", cert, key, nil); err != nil {
		fmt.Println(err, "failed to initialize a server")
	}
}
