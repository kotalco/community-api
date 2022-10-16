package ipfs_cluster_peer

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ipfsClusterPeerService struct{}

type IService interface {
	Get(name types.NamespacedName) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	Create(dto *ClusterPeerDto) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	Update(*ClusterPeerDto, *ipfsv1alpha1.ClusterPeer) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	List(namespace string) (*ipfsv1alpha1.ClusterPeerList, *restErrors.RestErr)
	Delete(node *ipfsv1alpha1.ClusterPeer) *restErrors.RestErr
	Count(namespace string) (*int, *restErrors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewIpfsClusterPeerService() IService {
	return ipfsClusterPeerService{}
}

// Get gets a single IPFS peer by name
func (service ipfsClusterPeerService) Get(namespacedName types.NamespacedName) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr) {
	peer := &ipfsv1alpha1.ClusterPeer{}

	if err := k8sClient.Get(context.Background(), namespacedName, peer); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("cluster peer by name %s doesn't exit", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get cluster peer by name %s", peer.Name))
	}

	return peer, nil
}

// Create creates IPFS peer from spec
func (service ipfsClusterPeerService) Create(dto *ClusterPeerDto) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr) {

	peer := &ipfsv1alpha1.ClusterPeer{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: ipfsv1alpha1.ClusterPeerSpec{
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if dto.PeerEndpoint != "" {
		peer.Spec.PeerEndpoint = dto.PeerEndpoint
	}

	if dto.Consensus != "" {
		peer.Spec.Consensus = ipfsv1alpha1.ConsensusAlgorithm(dto.Consensus)
	}

	if dto.ID != "" {
		peer.Spec.ID = dto.ID
	}

	if dto.PrivatekeySecretName != "" {
		peer.Spec.PrivateKeySecretName = dto.PrivatekeySecretName
	}

	if len(dto.TrustedPeers) != 0 {
		peer.Spec.TrustedPeers = dto.TrustedPeers
	}

	if len(dto.BootstrapPeers) != 0 {
		peer.Spec.BootstrapPeers = dto.BootstrapPeers
	}

	if dto.ClusterSecretName != "" {
		peer.Spec.ClusterSecretName = dto.ClusterSecretName
	}

	if dto.CPU != "" {
		peer.Spec.CPU = dto.CPU
	}
	if dto.CPULimit != "" {
		peer.Spec.CPULimit = dto.CPULimit
	}
	if dto.Memory != "" {
		peer.Spec.Memory = dto.Memory
	}
	if dto.MemoryLimit != "" {
		peer.Spec.MemoryLimit = dto.MemoryLimit
	}
	if dto.Storage != "" {
		peer.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8sClient.Create(context.Background(), peer); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("cluster peer by name %s already exits", peer.Name))
		}
		go logger.Error(service.Create, err)
		return nil, restErrors.NewInternalServerError("failed to create cluster peer")
	}

	return peer, nil
}

// Update updates IPFS peer by name from spec
func (service ipfsClusterPeerService) Update(dto *ClusterPeerDto, peer *ipfsv1alpha1.ClusterPeer) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr) {
	if dto.PeerEndpoint != "" {
		peer.Spec.PeerEndpoint = dto.PeerEndpoint
	}

	if len(dto.BootstrapPeers) != 0 {
		peer.Spec.BootstrapPeers = dto.BootstrapPeers
	}

	if dto.CPU != "" {
		peer.Spec.CPU = dto.CPU
	}
	if dto.CPULimit != "" {
		peer.Spec.CPULimit = dto.CPULimit
	}
	if dto.Memory != "" {
		peer.Spec.Memory = dto.Memory
	}
	if dto.MemoryLimit != "" {
		peer.Spec.MemoryLimit = dto.MemoryLimit
	}
	if dto.Storage != "" {
		peer.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8sClient.Update(context.Background(), peer); err != nil {
		go logger.Error(service.Update, err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't update cluster peer by name %s", peer.Name))
	}

	return peer, nil
}

// List returns all IPFS peers
func (service ipfsClusterPeerService) List(namespace string) (*ipfsv1alpha1.ClusterPeerList, *restErrors.RestErr) {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8sClient.List(context.Background(), peers, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return nil, restErrors.NewInternalServerError("failed to get all peers")
	}

	return peers, nil
}

// Count returns total number of IPFS peers
func (service ipfsClusterPeerService) Count(namespace string) (*int, *restErrors.RestErr) {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8sClient.List(context.Background(), peers, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.Count, err)
		return nil, restErrors.NewInternalServerError("failed to count all cluster perrs")
	}

	length := len(peers.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 IPFS peer by name
func (service ipfsClusterPeerService) Delete(peer *ipfsv1alpha1.ClusterPeer) *restErrors.RestErr {
	if err := k8sClient.Delete(context.Background(), peer); err != nil {
		go logger.Error(service.Delete, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delete cluster peer by name %s", peer.Name))
	}

	return nil
}
