// Package chainlink internal is the domain layer for the chainlink node
// uses the k8 client to CRUD the node
package chainlink

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type chainlinkService struct{}

type IService interface {
	Get(types.NamespacedName) (chainlinkv1alpha1.Node, restErrors.IRestErr)
	Create(ChainlinkDto) (chainlinkv1alpha1.Node, restErrors.IRestErr)
	Update(ChainlinkDto, *chainlinkv1alpha1.Node) restErrors.IRestErr
	List(namespace string) (chainlinkv1alpha1.NodeList, restErrors.IRestErr)
	Count(namespace string) (int, restErrors.IRestErr)
	Delete(*chainlinkv1alpha1.Node) restErrors.IRestErr
}

var (
	k8sClient = k8s.NewClientService()
)

func NewChainLinkService() IService {
	return chainlinkService{}
}

// Get returns a single chainlink node by name
func (service chainlinkService) Get(namespacedName types.NamespacedName) (node chainlinkv1alpha1.Node, restErr restErrors.IRestErr) {
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

// Create creates chainlink node from the given spec
func (service chainlinkService) Create(dto ChainlinkDto) (node chainlinkv1alpha1.Node, restErr restErrors.IRestErr) {

	node.ObjectMeta = dto.ObjectMetaFromMetadataDto()
	node.Spec = chainlinkv1alpha1.NodeSpec{
		EthereumChainId:            dto.EthereumChainId,
		LinkContractAddress:        dto.LinkContractAddress,
		EthereumWSEndpoint:         dto.EthereumWSEndpoint,
		DatabaseURL:                dto.DatabaseURL,
		KeystorePasswordSecretName: dto.KeystorePasswordSecretName,
		Image:                      dto.Image,
		APICredentials: chainlinkv1alpha1.APICredentials{
			Email:              dto.APICredentials.Email,
			PasswordSecretName: dto.APICredentials.PasswordSecretName,
		},
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

// Update updates a single chainlink node by name from spec
func (service chainlinkService) Update(dto ChainlinkDto, node *chainlinkv1alpha1.Node) (restErr restErrors.IRestErr) {

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
	if dto.API != nil {
		node.Spec.API = *dto.API
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

// List returns all chainlink nodes
func (service chainlinkService) List(namespace string) (list chainlinkv1alpha1.NodeList, restErr restErrors.IRestErr) {
	err := k8sClient.List(context.Background(), &list, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}

	return
}

// Count returns all nodes length
func (service chainlinkService) Count(namespace string) (count int, restErr restErrors.IRestErr) {
	nodes := &chainlinkv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		restErr = restErrors.NewInternalServerError("failed to count all nodes")
		return
	}

	return len(nodes.Items), nil
}

// Delete a single chainlink node by name
func (service chainlinkService) Delete(node *chainlinkv1alpha1.Node) (restErr restErrors.IRestErr) {
	err := k8sClient.Delete(context.Background(), node)

	if err != nil {
		go logger.Error(service.Delete, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}
	return
}
