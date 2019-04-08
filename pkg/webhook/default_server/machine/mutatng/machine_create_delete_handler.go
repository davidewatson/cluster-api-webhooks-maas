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
	"net/http"

	genericv1alpha1 "github.com/davidewatson/cluster-api-webhooks-maas/pkg/apis/generic/v1alpha1"
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
	// To use the client, you need to do the following:
	// - uncomment it
	// - import sigs.k8s.io/controller-runtime/pkg/client
	// - uncomment the InjectClient method at the bottom of this file.
	// Client  client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *MachineCreateDeleteHandler) mutatngMachineFn(ctx context.Context, obj *genericv1alpha1.Machine) (bool, string, error) {
	// TODO(user): implement your admission logic
	return true, "allowed to be admitted", nil
}

var _ admission.Handler = &MachineCreateDeleteHandler{}

// Handle handles admission requests.
func (h *MachineCreateDeleteHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &genericv1alpha1.Machine{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.mutatngMachineFn(ctx, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

//var _ inject.Client = &MachineCreateDeleteHandler{}
//
//// InjectClient injects the client into the MachineCreateDeleteHandler
//func (h *MachineCreateDeleteHandler) InjectClient(c client.Client) error {
//	h.Client = c
//	return nil
//}

var _ inject.Decoder = &MachineCreateDeleteHandler{}

// InjectDecoder injects the decoder into the MachineCreateDeleteHandler
func (h *MachineCreateDeleteHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
