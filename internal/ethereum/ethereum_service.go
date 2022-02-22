// Package ethereum internal is the domain layer for the ethereum node
// uses the k8 client to CRUD the node
package ethereum

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ethereumService struct{}

type ethereumServiceInterface interface {
	Get(name string) (*ethereumv1alpha1.Node, *errors.RestErr)
	Create(*EthereumDto) (*ethereumv1alpha1.Node, *errors.RestErr)
	Update(*EthereumDto, *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr)
	List() (*ethereumv1alpha1.NodeList, *errors.RestErr)
	Delete(node *ethereumv1alpha1.Node) *errors.RestErr
	Count() (*int, *errors.RestErr)
}

var (
	EthereumService ethereumServiceInterface
)

func init() { EthereumService = &ethereumService{} }

// Get returns a single ethereum node by name
func (service ethereumService) Get(name string) (*ethereumv1alpha1.Node, *errors.RestErr) {
	node := &ethereumv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}
	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", name))
		}
		go logger.Error("Error Getting ethereum Node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil

}

// Create creates ethereum node from the given spec
func (service ethereumService) Create(dto *EthereumDto) (*ethereumv1alpha1.Node, *errors.RestErr) {
	node := &ethereumv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
		},
		Spec: ethereumv1alpha1.NodeSpec{
			Network:                  dto.Network,
			Client:                   ethereumv1alpha1.EthereumClient(dto.Client),
			RPC:                      true,
			NodePrivateKeySecretName: dto.NodePrivateKeySecretName,
			Resources: sharedAPI.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Create(context.Background(), node)
	if err != nil {
		if k8errors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
		}
		go logger.Error("error creating ethereum node", err)
		return nil, errors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates a single ethereum node by name from spec
func (service ethereumService) Update(dto *EthereumDto, node *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr) {

	if dto.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(dto.Logging)
	}
	if dto.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = dto.NodePrivateKeySecretName
	}
	if dto.SyncMode != "" {
		node.Spec.SyncMode = ethereumv1alpha1.SynchronizationMode(dto.SyncMode)
	}
	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.Miner != nil {
		node.Spec.Miner = *dto.Miner
	}
	if node.Spec.Miner {
		if dto.Coinbase != "" {
			node.Spec.Coinbase = ethereumv1alpha1.EthereumAddress(dto.Coinbase)
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

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Update(context.Background(), node)
	if err != nil {
		go logger.Error("error updating ethereum node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all ethereum nodes
func (service ethereumService) List() (*ethereumv1alpha1.NodeList, *errors.RestErr) {
	nodes := &ethereumv1alpha1.NodeList{}

	err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default"))
	if err != nil {
		go logger.Error("Error listing ethereum nodes", err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns the length of ethereum nodes
func (service ethereumService) Count() (*int, *errors.RestErr) {
	nodes, err := service.List()
	if err != nil {
		return nil, err
	}

	length := len(nodes.Items)
	return &length, nil
}

// Delete a single ethereum node by name
func (service ethereumService) Delete(node *ethereumv1alpha1.Node) *errors.RestErr {
	err := k8s.Client().Delete(context.Background(), node)

	if err != nil {
		go logger.Error("Error deleting ethereum node", err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
