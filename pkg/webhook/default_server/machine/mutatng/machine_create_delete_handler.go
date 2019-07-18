/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mutatng

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cattlek8s/cluster-api-provider-generic/pkg/apis/generic/v1alpha1"
	"github.com/davidewatson/cluster-api-webhooks-maas/pkg/maas"
	//corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "mutatng-create-delete-machine"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &MachineCreateDeleteHandler{})
}

// MachineCreateDeleteHandler handles Machine
type MachineCreateDeleteHandler struct {
	// Client manipulates objects
	Client client.Client

	// Decoder decodes objects
	Decoder types.Decoder

	// MAASClient manages MAAS machines (e.g. allocate, deploy, release, etc.)
	MAASClient *maas.Client

	// TODO: Replace Client in this struct with this
	MachinesGetter client.Client
}

const (
	apiURLEnv     = "MAAS_API_URL"
	apiVersionEnv = "MAAS_API_VERSION"
	apiKeyEnv     = "MAAS_API_KEY"
)

// initMAASClient creates a MAAS Client and ClientSet for the webhook
// TODO: This should be done within the managers main() function.
func (h *MachineCreateDeleteHandler) initMAASClient() error {
	if h.MAASClient != nil {
		return nil
	}

	apiURL := os.Getenv(apiURLEnv)
	apiVersion := os.Getenv(apiVersionEnv)
	apiKey := os.Getenv(apiKeyEnv)
	maasClient, err := maas.New(&maas.ClientParams{ApiURL: apiURL, ApiVersion: apiVersion, ApiKey: apiKey})
	if err != nil {
		fmt.Printf("Unable to create MAAS client for machine controller: %v", err)
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	fmt.Printf("setting up client for manageri\n")
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("unable to set up client config: %v\n", err)
		os.Exit(1)
	}

	// Create ClientSet for webhook
	_, err = clientset.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("Unable to get clientset for config: %v\n", err)
		os.Exit(1)
	}

	h.MAASClient = &maasClient
	h.MachinesGetter = nil //cs.ClusterV1alpha1()

	return nil
}

func (h *MachineCreateDeleteHandler) mutatngMachineFn(ctx context.Context, obj *clusterv1.Machine) (bool, string, error) {
	if err := h.initMAASClient(); err != nil {
		panic("Unable to configure webhook")
	}

	response, err := h.MAASClient.Create(ctx, &v1alpha1.MachineCreateRequest{MachineID: obj.Name})
	if err != nil {
		return false, "webhook error prevents admission", err
	}

	obj.Spec.ProviderID = response.ProviderID
	// TODO: With the commented out code below we get this error:
	// > Error from server (InternalError): error when creating "config/samples/cluster_v1alpha1_machine.yaml": Internal error occurred: Internal error occurred: jsonpatch add operation does not apply: doc is missing path: "/status/addresses"
	//obj.Status.Addresses = make([]corev1.NodeAddress, 1)
	//obj.Status.Addresses[0] = corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: response.IPAddress}

	return true, "allowed to be admitted", nil
}

var _ admission.Handler = &MachineCreateDeleteHandler{}

// Handle handles admission requests.
func (h *MachineCreateDeleteHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &clusterv1.Machine{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := obj.DeepCopy()

	allowed, reason, err := h.mutatngMachineFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	if !allowed {
		return admission.ValidationResponse(allowed, reason)
	}

	res := admission.PatchResponse(obj, copy)
	fmt.Printf("res %v\n", res)
	return res
}

var _ inject.Client = &MachineCreateDeleteHandler{}

// InjectClient injects the client into the MachineCreateDeleteHandler
func (h *MachineCreateDeleteHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &MachineCreateDeleteHandler{}

// InjectDecoder injects the decoder into the MachineCreateDeleteHandler
func (h *MachineCreateDeleteHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
