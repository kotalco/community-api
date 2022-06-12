// Package chainlink internal is the domain layer for the chainlink node
// uses the k8 client to CRUD the node
package chainlink

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type chainlinkService struct{}

type IService interface {
	Get(types.NamespacedName) (*chainlinkv1alpha1.Node, *errors.RestErr)
	Create(*ChainlinkDto) (*chainlinkv1alpha1.Node, *errors.RestErr)
	Update(*ChainlinkDto, *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr)
	List(namespace string) (*chainlinkv1alpha1.NodeList, *errors.RestErr)
	Count(namespace string) (*int, *errors.RestErr)
	Delete(node *chainlinkv1alpha1.Node) *errors.RestErr
}

var (
	k8sClient = k8s.NewClientService()
)

func NewChainLinkService() IService {
	return chainlinkService{}
}

// Get returns a single chainlink node by name
func (service chainlinkService) Get(namespacedName types.NamespacedName) (*chainlinkv1alpha1.Node, *errors.RestErr) {
	node := &chainlinkv1alpha1.Node{}
	if err := k8sClient.Get(context.Background(), namespacedName, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
	}

	return node, nil

}

// Create creates chainlink node from the given spec
func (service chainlinkService) Create(dto *ChainlinkDto) (*chainlinkv1alpha1.Node, *errors.RestErr) {
	node := &chainlinkv1alpha1.Node{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: chainlinkv1alpha1.NodeSpec{
			EthereumChainId:            dto.EthereumChainId,
			LinkContractAddress:        dto.LinkContractAddress,
			EthereumWSEndpoint:         dto.EthereumWSEndpoint,
			DatabaseURL:                dto.DatabaseURL,
			KeystorePasswordSecretName: dto.KeystorePasswordSecretName,
			APICredentials: chainlinkv1alpha1.APICredentials{
				Email:              dto.APICredentials.Email,
				PasswordSecretName: dto.APICredentials.PasswordSecretName,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8sClient.Create(context.Background(), node)
	if err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
		}
		go logger.Error(service.Create, err)
		return nil, errors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates a single chainlink node by name from spec
func (service chainlinkService) Update(dto *ChainlinkDto, node *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr) {

	if dto.EthereumWSEndpoint != "" {
		node.Spec.EthereumWSEndpoint = dto.EthereumWSEndpoint
	}

	if dto.DatabaseURL != "" {
		node.Spec.DatabaseURL = dto.DatabaseURL
	}

	if len(dto.EthereumHTTPEndpoints) != 0 {
		node.Spec.EthereumHTTPEndpoints = dto.EthereumHTTPEndpoints
	}

	if dto.KeystorePasswordSecretName != "" {
		node.Spec.KeystorePasswordSecretName = dto.KeystorePasswordSecretName
	}

	if dto.APICredentials != nil {
		node.Spec.APICredentials.Email = dto.APICredentials.Email
		node.Spec.APICredentials.PasswordSecretName = dto.APICredentials.PasswordSecretName
	}

	if len(dto.CORSDomains) != 0 {
		node.Spec.CORSDomains = dto.CORSDomains
	}

	if dto.CertSecretName != "" {
		node.Spec.CertSecretName = dto.CertSecretName
	}

	if dto.TLSPort != 0 {
		node.Spec.TLSPort = dto.TLSPort
	}

	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.APIPort != 0 {
		node.Spec.APIPort = dto.APIPort
	}

	if dto.SecureCookies != nil {
		node.Spec.SecureCookies = *dto.SecureCookies
	}

	if dto.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(dto.Logging)
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

	err := k8sClient.Update(context.Background(), node)
	if err != nil {
		go logger.Error(service.Update, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all chainlink nodes
func (service chainlinkService) List(namespace string) (*chainlinkv1alpha1.NodeList, *errors.RestErr) {
	nodes := &chainlinkv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns all nodes length
func (service chainlinkService) Count(namespace string) (*int, *errors.RestErr) {
	nodes := &chainlinkv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	length := len(nodes.Items)
	return &length, nil
}

// Delete a single chainlink node by name
func (service chainlinkService) Delete(node *chainlinkv1alpha1.Node) *errors.RestErr {
	err := k8sClient.Delete(context.Background(), node)

	if err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
