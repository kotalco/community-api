package filecoin

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type filecoinService struct{}

type IService interface {
	Get(types.NamespacedName) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	Create(dto *FilecoinDto) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	Update(*FilecoinDto, *filecoinv1alpha1.Node) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	List(namespace string) (*filecoinv1alpha1.NodeList, *restErrors.RestErr)
	Delete(node *filecoinv1alpha1.Node) *restErrors.RestErr
	Count(namespace string) (*int, *restErrors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewFilecoinService() IService {
	return filecoinService{}
}

// Get gets a single filecoin node by name
func (service filecoinService) Get(namespacedName types.NamespacedName) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	node := &filecoinv1alpha1.Node{}

	if err := k8sClient.Get(context.Background(), namespacedName, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exit", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
	}

	return node, nil
}

// Create creates filecoin node from spec
func (service filecoinService) Create(dto *FilecoinDto) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	node := &filecoinv1alpha1.Node{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: filecoinv1alpha1.NodeSpec{
			Network: filecoinv1alpha1.FilecoinNetwork(dto.Network),
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8sClient.Create(context.Background(), node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("node by name %s already exits", dto))
		}
		go logger.Error(service.Create, err)
		return nil, restErrors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates filecoin node by name from spec
func (service filecoinService) Update(dto *FilecoinDto, node *filecoinv1alpha1.Node) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	if dto.API != nil {
		node.Spec.API = *dto.API
	}

	if dto.APIPort != 0 {
		node.Spec.APIPort = dto.APIPort
	}

	if dto.APIHost != "" {
		node.Spec.APIHost = dto.APIHost
	}

	if dto.APIRequestTimeout != 0 {
		node.Spec.APIRequestTimeout = dto.APIRequestTimeout
	}

	if dto.DisableMetadataLog != nil {
		node.Spec.DisableMetadataLog = *dto.DisableMetadataLog
	}

	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.P2PHost != "" {
		node.Spec.P2PHost = dto.P2PHost
	}

	if dto.IPFSPeerEndpoint != "" {
		node.Spec.IPFSPeerEndpoint = dto.IPFSPeerEndpoint
	}

	if dto.IPFSOnlineMode != nil {
		node.Spec.IPFSOnlineMode = *dto.IPFSOnlineMode
	}

	if dto.IPFSForRetrieval != nil {
		node.Spec.IPFSForRetrieval = *dto.IPFSForRetrieval
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
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all filecoin nodes
func (service filecoinService) List(namespace string) (*filecoinv1alpha1.NodeList, *restErrors.RestErr) {
	nodes := &filecoinv1alpha1.NodeList{}
	if err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return nil, restErrors.NewInternalServerError("failed to get all nodes")
	}
	return nodes, nil
}

// Count returns total number of filecoin nodes
func (service filecoinService) Count(namespace string) (*int, *restErrors.RestErr) {

	nodes := &filecoinv1alpha1.NodeList{}
	if err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.Count, err)
		return nil, restErrors.NewInternalServerError("failed to count filecoin nodes")
	}

	length := len(nodes.Items)
	return &length, nil
}

// Delete deletes ethereum 2.0 filecoin node by name
func (service filecoinService) Delete(node *filecoinv1alpha1.Node) *restErrors.RestErr {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delte node by name %s", node.Name))
	}
	return nil
}
