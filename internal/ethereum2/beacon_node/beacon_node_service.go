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
	corev1 "k8s.io/api/core/v1"
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
	client := ethereum2v1alpha1.Ethereum2Client(dto.Client)

	beaconnode := &ethereum2v1alpha1.BeaconNode{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: ethereum2v1alpha1.BeaconNodeSpec{
<<<<<<< HEAD
			Network:                 dto.Network,
			Client:                  client,
			RPC:                     client == ethereum2v1alpha1.PrysmClient,
			ExecutionEngineEndpoint: dto.ExecutionEngineEndpoint,
			JWTSecretName:           dto.JWTSecretName,
=======
			Network: dto.Network,
			Client:  client,
			RPC:     client == ethereum2v1alpha1.PrysmClient,
			Image:   dto.Image,
>>>>>>> 851c69f (feat: all procols can set or update image version (closing #47))
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
<<<<<<< HEAD
	if dto.ExecutionEngineEndpoint != "" {
		node.Spec.ExecutionEngineEndpoint = dto.ExecutionEngineEndpoint
	}
	if dto.JWTSecretName != "" {
		node.Spec.JWTSecretName = dto.JWTSecretName
=======
	if *dto.Image != "" {
		node.Spec.Image = dto.Image
>>>>>>> 851c69f (feat: all procols can set or update image version (closing #47))
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
			return nil, errors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
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
