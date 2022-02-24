package ipfs_cluster_peer

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ipfsClusterPeerService struct{}

type ipfsClusterPeerServiceInterface interface {
	Get(name string) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	Create(dto *ClusterPeerDto) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	Update(*ClusterPeerDto, *ipfsv1alpha1.ClusterPeer) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr)
	List() (*ipfsv1alpha1.ClusterPeerList, *restErrors.RestErr)
	Delete(node *ipfsv1alpha1.ClusterPeer) *restErrors.RestErr
	Count() (*int, *restErrors.RestErr)
}

var (
	IpfsClusterPeerService ipfsClusterPeerServiceInterface
)

func init() { IpfsClusterPeerService = &ipfsClusterPeerService{} }

// Get gets a single IPFS peer by name
func (service ipfsClusterPeerService) Get(name string) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr) {
	peer := &ipfsv1alpha1.ClusterPeer{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, peer); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("cluster peer by name %s doesn't exit", name))
		}
		go logger.Error("ERROR_IN_GET_SERVICE", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get cluster peer by name %s", peer.Name))
	}

	return peer, nil
}

// Create creates IPFS peer from spec
func (service ipfsClusterPeerService) Create(dto *ClusterPeerDto) (*ipfsv1alpha1.ClusterPeer, *restErrors.RestErr) {

	peer := &ipfsv1alpha1.ClusterPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:              dto.Name,
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
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

	if err := k8s.Client().Create(context.Background(), peer); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("cluster peer by name %s already exits", peer.Name))
		}
		go logger.Error("ERROR_IN_ CREATE_IPFS_CLUSTER_PEER", err)
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

	if err := k8s.Client().Update(context.Background(), peer); err != nil {
		go logger.Error("ERROR_IN_UPDATE_IPFS_CLUSTER_PEER", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't update cluster peer by name %s", peer.Name))
	}

	return peer, nil
}

// List returns all IPFS peers
func (service ipfsClusterPeerService) List() (*ipfsv1alpha1.ClusterPeerList, *restErrors.RestErr) {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8s.Client().List(context.Background(), peers, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_LIST_IPFS_CLUSTER_PEER_SERVICE", err)
		return nil, restErrors.NewInternalServerError("failed to get all peers")
	}

	return peers, nil
}

// Count returns total number of IPFS peers
func (service ipfsClusterPeerService) Count() (*int, *restErrors.RestErr) {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8s.Client().List(context.Background(), peers, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_COUNT_CLUSTER_PEER_SERVICE", err)
		return nil, restErrors.NewInternalServerError("failed to count all cluster perrs")
	}

	length := len(peers.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 IPFS peer by name
func (service ipfsClusterPeerService) Delete(peer *ipfsv1alpha1.ClusterPeer) *restErrors.RestErr {
	if err := k8s.Client().Delete(context.Background(), peer); err != nil {
		go logger.Error("ERROR_IN_DELETE_CLUSTER_PEER_SERVICE", err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delete cluster peer by name %s", peer.Name))
	}

	return nil
}
