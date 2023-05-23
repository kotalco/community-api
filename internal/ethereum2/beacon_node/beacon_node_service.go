// Package beacon_node internal is the domain layer for the ethereum2 beaconnode
// uses the k8 client to CRUD the node
package beacon_node

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type beaconNodeService struct{}

type IService interface {
	Get(types.NamespacedName) (ethereum2v1alpha1.BeaconNode, restErrors.IRestErr)
	Create(dto BeaconNodeDto) (ethereum2v1alpha1.BeaconNode, restErrors.IRestErr)
	Update(BeaconNodeDto, *ethereum2v1alpha1.BeaconNode) restErrors.IRestErr
	List(namespace string) (ethereum2v1alpha1.BeaconNodeList, restErrors.IRestErr)
	Delete(*ethereum2v1alpha1.BeaconNode) restErrors.IRestErr
	Count(namespace string) (int, restErrors.IRestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewBeaconNodeService() IService {
	return beaconNodeService{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (service beaconNodeService) Get(namespacedNamed types.NamespacedName) (node ethereum2v1alpha1.BeaconNode, restErr restErrors.IRestErr) {
	if err := k8sClient.Get(context.Background(), namespacedNamed, &node); err != nil {
		if apiErrors.IsNotFound(err) {
			restErr = restErrors.NewNotFoundError(fmt.Sprintf("beacon node by name %s doesn't exist", namespacedNamed.Name))
			return
		}
		go logger.Error(service.Get, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't get beacon node by name %s", namespacedNamed.Name))
		return
	}

	return
}

// Create creates ethereum 2.0 beacon node from spec
func (service beaconNodeService) Create(dto BeaconNodeDto) (node ethereum2v1alpha1.BeaconNode, restErr restErrors.IRestErr) {
	client := ethereum2v1alpha1.Ethereum2Client(dto.Client)

	node.ObjectMeta = dto.ObjectMetaFromMetadataDto()
	node.Spec = ethereum2v1alpha1.BeaconNodeSpec{
		Network:                 dto.Network,
		Client:                  client,
		RPC:                     client == ethereum2v1alpha1.PrysmClient,
		ExecutionEngineEndpoint: dto.ExecutionEngineEndpoint,
		JWTSecretName:           dto.JWTSecretName,
		Image:                   dto.Image,
		REST:                    true,
		Resources: sharedAPIs.Resources{
			StorageClass: dto.StorageClass,
		},
	}

	k8s.DefaultResources(&node.Spec.Resources)

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8sClient.Create(context.Background(), &node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			restErr = restErrors.NewBadRequestError(fmt.Sprintf("beacon node by name %s already exist", dto.Name))
			return
		}
		go logger.Error(service.Create, err)
		restErr = restErrors.NewInternalServerError("failed to create beacon node")
		return
	}

	return
}

// Update updates ethereum 2.0 beacon node by name from spec
func (service beaconNodeService) Update(dto BeaconNodeDto, node *ethereum2v1alpha1.BeaconNode) (restErr restErrors.IRestErr) {
	if dto.REST != nil {
		rest := *dto.REST
		if rest {
			if dto.RESTPort != 0 {
				node.Spec.RESTPort = dto.RESTPort
			}
		}
		node.Spec.REST = rest
	}

	if dto.RPC != nil {
		rpc := *dto.RPC
		if rpc {
			if dto.RPCPort != 0 {
				node.Spec.RPCPort = dto.RPCPort
			}
		}
		node.Spec.RPC = rpc
	}

	if dto.GRPC != nil {
		grpc := *dto.GRPC
		if grpc {
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
	if dto.ExecutionEngineEndpoint != "" {
		node.Spec.ExecutionEngineEndpoint = dto.ExecutionEngineEndpoint
	}
	if dto.JWTSecretName != "" {
		node.Spec.JWTSecretName = dto.JWTSecretName
	}
	if dto.Image != "" {
		node.Spec.Image = dto.Image
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	pod := &corev1.Pod{}
	podIsPending := false
	if dto.CPU != "" || dto.Memory != "" {
		key := types.NamespacedName{
			Namespace: node.Namespace,
			Name:      fmt.Sprintf("%s-0", node.Name),
		}
		err := k8sClient.Get(context.Background(), key, pod)
		if apiErrors.IsNotFound(err) {
			go logger.Error(service.Update, err)
			restErr = restErrors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
			return
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
		return
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
			return
		}
	}

	return
}

// List returns all ethereum 2.0 beacon nodes
func (service beaconNodeService) List(namespace string) (list ethereum2v1alpha1.BeaconNodeList, restErr restErrors.IRestErr) {
	if err := k8sClient.List(context.Background(), &list, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		restErr = restErrors.NewInternalServerError("failed to get all beacon nodes")
		return
	}

	return
}

// Count returns total number of beacon nodes
func (service beaconNodeService) Count(namespace string) (count int, restErr restErrors.IRestErr) {
	nodes, err := service.List(namespace)
	if err != nil {
		restErr = restErrors.NewInternalServerError("failed to count all nodes")
		return
	}

	return len(nodes.Items), nil
}

// Delete deletes ethereum 2.0 beacon node by name
func (service beaconNodeService) Delete(node *ethereum2v1alpha1.BeaconNode) (restErr restErrors.IRestErr) {
	err := k8sClient.Delete(context.Background(), node)

	if err != nil {
		go logger.Error(service.Delete, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
		return
	}

	return
}
