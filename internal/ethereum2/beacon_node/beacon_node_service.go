// Package beacon_node internal is the domain layer for the ethereum2 beaconnode
// uses the k8 client to CRUD the node
package beacon_node

import (
	"context"
	"fmt"
	"github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type beaconNodeService struct{}

type IService interface {
	Get(types.NamespacedName) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr)
	Create(dto *BeaconNodeDto) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr)
	Update(*BeaconNodeDto, *ethereum2v1alpha1.BeaconNode) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr)
	List(namespace string) (*ethereum2v1alpha1.BeaconNodeList, *errors.RestErr)
	Delete(node *ethereum2v1alpha1.BeaconNode) *errors.RestErr
	Count(namespace string) (*int, *errors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewBeaconNodeService() IService {
	return beaconNodeService{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (service beaconNodeService) Get(namespacedNamed types.NamespacedName) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr) {
	node := &ethereum2v1alpha1.BeaconNode{}

	if err := k8sClient.Get(context.Background(), namespacedNamed, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("beacon node by name %s doesn't exist", namespacedNamed.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get beacon node by name %s", namespacedNamed.Name))
	}

	return node, nil
}

// Create creates ethereum 2.0 beacon node from spec
func (service beaconNodeService) Create(dto *BeaconNodeDto) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr) {

	var endpoints []string
	if dto.Eth1Endpoints != nil {
		endpoints = *dto.Eth1Endpoints
	}

	client := ethereum2v1alpha1.Ethereum2Client(dto.Client)

	beaconnode := &ethereum2v1alpha1.BeaconNode{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: ethereum2v1alpha1.BeaconNodeSpec{
			Network:       dto.Network,
			Client:        client,
			Eth1Endpoints: endpoints,
			RPC:           client == ethereum2v1alpha1.PrysmClient,
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	k8s.DefaultResources(&beaconnode.Spec.Resources)

	if os.Getenv("MOCK") == "true" {
		beaconnode.Default()
	}

	if err := k8sClient.Create(context.Background(), beaconnode); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("beacon node by name %s already exist", dto.Name))
		}
		go logger.Error(service.Create, err)
		return nil, errors.NewInternalServerError("failed to create beacon node")
	}

	return beaconnode, nil
}

// Update updates ethereum 2.0 beacon node by name from spec
func (service beaconNodeService) Update(dto *BeaconNodeDto, node *ethereum2v1alpha1.BeaconNode) (*ethereum2v1alpha1.BeaconNode, *errors.RestErr) {
	endpoints := dto.Eth1Endpoints
	if endpoints != nil {
		// all clients can clear ethereum endpoints
		// prysm can clear ethereum endpoints only if network is mainnet
		if node.Spec.Client == ethereum2v1alpha1.PrysmClient && node.Spec.Network != "mainnet" && len(*endpoints) == 0 {
			// do nothing
		} else {
			node.Spec.Eth1Endpoints = *endpoints
		}
	}

	if dto.REST != nil {
		rest := *dto.REST
		if rest {
			if dto.RESTHost != "" {
				node.Spec.RESTHost = dto.RESTHost
			}
			if dto.RESTPort != 0 {
				node.Spec.RESTPort = dto.RESTPort
			}
		}
		node.Spec.REST = rest
	}

	if dto.RPC != nil {
		rpc := *dto.RPC
		if rpc {
			if dto.RPCHost != "" {
				node.Spec.RPCHost = dto.RPCHost
			}
			if dto.RPCPort != 0 {
				node.Spec.RPCPort = dto.RPCPort
			}
		}
		node.Spec.RPC = rpc
	}

	if dto.GRPC != nil {
		grpc := *dto.GRPC
		if grpc {
			if dto.GRPCHost != "" {
				node.Spec.GRPCHost = dto.GRPCHost
			}
			if dto.GRPCPort != 0 {
				node.Spec.GRPCPort = dto.GRPCPort
			}
		}
		node.Spec.GRPC = grpc
	}

	if dto.CPU != "" {
		node.Spec.CPU = dto.CPU
	}
	if dto.CPULimit != "" {
		node.Spec.CPULimit = dto.CPULimit
	}
	if dto.Memory != "" {
		node.Spec.Memory = dto.Memory
	}
	if dto.MemoryLimit != "" {
		node.Spec.MemoryLimit = dto.MemoryLimit
	}
	if dto.Storage != "" {
		node.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name  %s", node.Name))
	}

	if k8s.CheckDeploymentResourcesChanged(&dto.Resources) {
		err := k8s.DeployReconciliation(node.Name, node.Namespace)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

// List returns all ethereum 2.0 beacon nodes
func (service beaconNodeService) List(namespace string) (*ethereum2v1alpha1.BeaconNodeList, *errors.RestErr) {
	nodes := &ethereum2v1alpha1.BeaconNodeList{}

	if err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all beacon nodes")
	}

	return nodes, nil
}

// Count returns total number of beacon nodes
func (service beaconNodeService) Count(namespace string) (*int, *errors.RestErr) {
	nodes, err := service.List(namespace)
	if err != nil {
		return nil, err
	}

	length := len(nodes.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 beacon node by name
func (service beaconNodeService) Delete(node *ethereum2v1alpha1.BeaconNode) *errors.RestErr {
	err := k8sClient.Delete(context.Background(), node)

	if err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
