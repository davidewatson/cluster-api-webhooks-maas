package maas

import (
	"context"
	"fmt"

	"github.com/cattlek8s/cluster-api-provider-generic/pkg/apis/generic/v1alpha1"
	"k8s.io/klog"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterv1client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	"github.com/juju/gomaasapi"
)

type Client struct {
	Controller gomaasapi.Controller
}

type ClientParams struct {
	ApiURL         string
	ApiVersion     string
	ApiKey         string
	V1Alpha1Client clusterv1client.ClusterV1alpha1Interface
}

func New(params *ClientParams) (Client, error) {
	controller, err := gomaasapi.NewController(gomaasapi.ControllerArgs{
		BaseURL: params.ApiURL,
		APIKey:  params.ApiKey})
	if err != nil {
		return Client{}, fmt.Errorf("error creating controller with version: %v", err)
	}

	return Client{Controller: controller}, nil
}

// Create creates a machine
func (c Client) Create(ctx context.Context, request *v1alpha1.MachineCreateRequest) (*v1alpha1.MachineCreateResponse, error) {
	klog.Infof("Creating machine %s", request.MachineID)

	// Allocate MAAS machine
	allocateArgs := gomaasapi.AllocateMachineArgs{Tags: []string{}}
	m, _, err := c.Controller.AllocateMachine(allocateArgs)
	if err != nil {
		klog.Errorf("Create failed to allocate machine %s: %v", request.MachineID, err)
		return nil, fmt.Errorf("error allocating machine %s: %v", request.MachineID, err)
	}

	// Deploy MAAS machine
	startArgs := gomaasapi.StartArgs{
		DistroSeries: "ubuntu-18.04-cnct-k8s-master", // TODO: Must depend on k8s version?
	}
	err = m.Start(startArgs)

	// Release if there are any errors...
	if err != nil {
		klog.Errorf("Create failed to deploy machine %s: %v", request.MachineID, err)

		err := c.Delete(ctx, &DeleteRequest{MachineID: request.MachineID, ProviderID: m.SystemID()})
		if err != nil {
			klog.Errorf("Create failed to release machine %s: %v", request.MachineID, err)
		}

		return nil, err
	}

	providerID := m.SystemID()
	ipAddresses := m.IPAddresses()

	if len(ipAddresses) < 1 {
		klog.Errorf("Create failed to deploy machine %s: %v", request.MachineID, err)

		err := c.Delete(ctx, &DeleteRequest{MachineID: request.MachineID,
			ProviderID: providerID})
		if err != nil {
			klog.Errorf("Create failed to release machine %s: %v", request.MachineID, err)
		}

		return nil, err
	}

	// Success
	klog.Infof("Created machine %s with IP address %s", providerID, ipAddresses[0])

	return &v1alpha1.MachineCreateResponse{
		ProviderID: &providerID,
		IPAddress:  ipAddresses[0],
	}, nil
}

type DeleteRequest struct {
	// MachineID is the unique value passed in CreateRequest.
	MachineID string
	// SystemID is the unique value passed in CreateResponse.
	ProviderID string
}

type DeleteResponse struct {
}

// Delete deletes a machine
func (c Client) Delete(ctx context.Context, request *DeleteRequest) error {
	if request.ProviderID == "" {
		klog.Warningf("can not delete  machine %s, providerID not set", request.MachineID)
		return fmt.Errorf("machine %s has not been created", request.MachineID)
	}

	// Release MAAS machine
	releaseArgs := gomaasapi.ReleaseMachinesArgs{SystemIDs: []string{request.ProviderID}}
	if err := c.Controller.ReleaseMachines(releaseArgs); err != nil {
		klog.Warningf("error releasing machine %s (%s): %v", request.MachineID, request.ProviderID, err)
		return nil
	}

	return nil
}

// Update updates a machine
func (c Client) Update(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

// Exists test for the existence of a machine
func (c Client) Exist(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	// Get list of machines with tag
	machineArgs := gomaasapi.MachinesArgs{SystemIDs: []string{*machine.Spec.ProviderID}}
	machines, err := c.Controller.Machines(machineArgs)
	if err != nil {
		return false, fmt.Errorf("error listing machines: %v", err)
	}
	if len(machines) != 1 {
		return false, fmt.Errorf("expected 1 machine (%s), found %d", machine.Name, len(machines))
	}

	return true, nil
}
