// Package ethereum internal is the domain layer for the ethereum node
// uses the k8 client to CRUD the node
package ethereum

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ethereumService struct{}

type IService interface {
	Get(types.NamespacedName) (ethereumv1alpha1.Node, restErrors.IRestErr)
	Create(EthereumDto) (ethereumv1alpha1.Node, restErrors.IRestErr)
	Update(EthereumDto, *ethereumv1alpha1.Node) restErrors.IRestErr
	List(namespace string) (ethereumv1alpha1.NodeList, restErrors.IRestErr)
	Delete(*ethereumv1alpha1.Node) restErrors.IRestErr
	Count(namespace string) (int, restErrors.IRestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewEthereumService() IService {
	return ethereumService{}
}

// Get returns a single ethereum node by name
func (service ethereumService) Get(namespacedName types.NamespacedName) (node ethereumv1alpha1.Node, restErr restErrors.IRestErr) {
	if err := k8sClient.Get(context.Background(), namespacedName, &node); err != nil {
		if apiErrors.IsNotFound(err) {
			restErr = restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", namespacedName.Name))
			return
		}
		go logger.Error(service.Get, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
		return
	}

	return
}

// Create creates ethereum node from the given spec
func (service ethereumService) Create(dto EthereumDto) (node ethereumv1alpha1.Node, restErr restErrors.IRestErr) {
	node.ObjectMeta = dto.ObjectMetaFromMetadataDto()
	node.Spec = ethereumv1alpha1.NodeSpec{
		Network: dto.Network,
		Client:  ethereumv1alpha1.EthereumClient(dto.Client),
		RPC:     true,
		Image:   dto.Image,
	}

	k8s.DefaultResources(&node.Spec.Resources)

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8sClient.Create(context.Background(), &node)
	if err != nil {
		if apiErrors.IsAlreadyExists(err) {
			restErr = restErrors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
			return
		}
		go logger.Error(service.Create, err)
		restErr = restErrors.NewInternalServerError("failed to create node")
		return
	}

	return
}

// Update updates a single ethereum node by name from spec
func (service ethereumService) Update(dto EthereumDto, node *ethereumv1alpha1.Node) (restErr restErrors.IRestErr) {

	if dto.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(dto.Logging)
	}
	if dto.NodePrivateKeySecretName != nil {
		node.Spec.NodePrivateKeySecretName = *dto.NodePrivateKeySecretName
	}
	if dto.SyncMode != nil {
		node.Spec.SyncMode = *dto.SyncMode
	}
	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.Miner != nil {
		node.Spec.Miner = *dto.Miner
	}
	if node.Spec.Miner {
		if dto.Coinbase != "" {
			node.Spec.Coinbase = sharedAPI.EthereumAddress(dto.Coinbase)
		}
		if dto.Import != nil {
			node.Spec.Import = &ethereumv1alpha1.ImportedAccount{
				PrivateKeySecretName: dto.Import.PrivateKeySecretName,
				PasswordSecretName:   dto.Import.PasswordSecretName,
			}
		}
	}

	if dto.RPC != nil {
		node.Spec.RPC = *dto.RPC
	}
	if node.Spec.RPC {
		if len(dto.RPCAPI) != 0 {
			rpcAPI := []ethereumv1alpha1.API{}
			for _, api := range dto.RPCAPI {
				rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.RPCAPI = rpcAPI
		}
		if dto.RPCPort != 0 {
			node.Spec.RPCPort = dto.RPCPort
		}
	}

	if dto.WS != nil {
		node.Spec.WS = *dto.WS
	}
	if node.Spec.WS {
		if len(dto.WSAPI) != 0 {
			wsAPI := []ethereumv1alpha1.API{}
			for _, api := range dto.WSAPI {
				wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.WSAPI = wsAPI
		}
		if dto.WSPort != 0 {
			node.Spec.WSPort = dto.WSPort
		}
	}

	if dto.GraphQL != nil {
		node.Spec.GraphQL = *dto.GraphQL
	}
	if node.Spec.GraphQL {
		if dto.GraphQLPort != 0 {
			node.Spec.GraphQLPort = dto.GraphQLPort
		}
	}

	if len(dto.Hosts) != 0 {
		node.Spec.Hosts = dto.Hosts
	}

	if len(dto.CORSDomains) != 0 {
		node.Spec.CORSDomains = dto.CORSDomains
	}

	var bootnodes, staticNodes []ethereumv1alpha1.Enode

	if dto.Bootnodes != nil {
		for _, bootnode := range *dto.Bootnodes {
			bootnodes = append(bootnodes, ethereumv1alpha1.Enode(bootnode))
		}
	}
	node.Spec.Bootnodes = bootnodes

	if dto.StaticNodes != nil {
		for _, staticNode := range *dto.StaticNodes {
			staticNodes = append(staticNodes, ethereumv1alpha1.Enode(staticNode))
		}
	}
	node.Spec.StaticNodes = staticNodes

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
	if dto.Engine != nil {
		node.Spec.Engine = *dto.Engine
	}
	if node.Spec.Engine {
		if dto.JWTSecretName != "" {
			node.Spec.JWTSecretName = dto.JWTSecretName
		}
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

// List returns all ethereum nodes
func (service ethereumService) List(namespace string) (list ethereumv1alpha1.NodeList, restErr restErrors.IRestErr) {
	err := k8sClient.List(context.Background(), &list, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}
	return
}

// Count returns the length of ethereum nodes
func (service ethereumService) Count(namespace string) (count int, restErr restErrors.IRestErr) {
	nodes, err := service.List(namespace)
	if err != nil {
		restErr = restErrors.NewInternalServerError("failed to count all nodes")
		return
	}

	return len(nodes.Items), nil
}

// Delete a single ethereum node by name
func (service ethereumService) Delete(node *ethereumv1alpha1.Node) (restErr restErrors.IRestErr) {
	err := k8sClient.Delete(context.Background(), node)

	if err != nil {
		go logger.Error(service.Delete, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
		return
	}

	return
}
