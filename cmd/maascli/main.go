package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/cattlek8s/cluster-api-provider-generic/pkg/apis/generic/v1alpha1"
	"github.com/davidewatson/cluster-api-webhooks-maas/pkg/maas"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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

	_, err = client.Create(context.TODO(), &v1alpha1.MachineCreateRequest{MachineID: fmt.Sprintf("%s-%s", "cluster1", "machine1")})
	if err != nil {
		klog.Fatalf("Failed to create Machine: %v\n", err)
	}
}
