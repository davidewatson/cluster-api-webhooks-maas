package main

import (
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/davidewatson/cluster-api-webhooks-maas/pkg/maas"
)

const (
	apiUrlKey     = "api_url"
	apiVersionKey = "api_version"
	apiKeyKey     = "api_key"
)

// init configures input and output.
func init() {
	viper.SetEnvPrefix("maas")
	viper.BindEnv(apiUrlKey)
	viper.BindEnv(apiVersionKey)
	viper.BindEnv(apiKeyKey)
	viper.AutomaticEnv()
}

// Prints out connection information and verifies it can be used to connect...
func main() {
	klog.InitFlags(nil)

	apiUrl := viper.GetString(apiUrlKey)
	apiVersion := viper.GetString(apiVersionKey)
	apiKey := viper.GetString(apiKeyKey)

	fmt.Printf("%s: %s\n%s: %s\n%s: %s\n", apiUrlKey, apiUrl, apiVersionKey, apiVersion, apiKeyKey, apiKey)

	_, err := maas.New(apiUrl, apiVersion, apiKey)
	if err != nil {
		klog.Fatalf("failed to create MAAS client: %v\n", err)
	}
}
