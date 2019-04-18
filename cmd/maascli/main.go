package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/klog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/davidewatson/cluster-api-webhooks-maas/pkg/maas"
)

const (
	apiURLKey     = "api_url"
	apiVersionKey = "api_version"
	apiKeyKey     = "api_key"
)

// init configures input and output.
func init() {
	viper.SetEnvPrefix("maas")
	viper.BindEnv(apiURLKey)
	viper.BindEnv(apiVersionKey)
	viper.BindEnv(apiKeyKey)
	viper.AutomaticEnv()
}

// Prints out connection information and verifies it can be used to connect...
func main() {
	klog.InitFlags(nil)

	cfg := config.GetConfigOrDie()
	cs, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Failed to create client from configuration: %v", err)
	}

	apiURL := viper.GetString(apiURLKey)
	apiVersion := viper.GetString(apiVersionKey)
	apiKey := viper.GetString(apiKeyKey)

	fmt.Printf("%s: %s\n%s: %s\n%s: %s\n", apiURLKey, apiURL, apiVersionKey, apiVersion, apiKeyKey, apiKey)

	client, err := maas.New(&maas.ClientParams{
		ApiURL:         apiURL,
		ApiVersion:     apiVersion,
		ApiKey:         apiKey,
		V1Alpha1Client: cs.ClusterV1alpha1()})
	if err != nil {
		klog.Fatalf("Failed to create MAAS client: %v\n", err)
	}

	err = client.Create(context.TODO(), nil, &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: corev1.NamespaceDefault}})
	if err != nil {
		klog.Fatalf("Failed to create Machine: %v\n", err)
	}
}
