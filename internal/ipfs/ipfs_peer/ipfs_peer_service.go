package ipfs_peer

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ipfsPeerService struct{}

type IService interface {
	Get(name types.NamespacedName) (ipfsv1alpha1.Peer, restErrors.IRestErr)
	Create(PeerDto) (ipfsv1alpha1.Peer, restErrors.IRestErr)
	Update(PeerDto, *ipfsv1alpha1.Peer) restErrors.IRestErr
	List(namespace string) (ipfsv1alpha1.PeerList, restErrors.IRestErr)
	Delete(*ipfsv1alpha1.Peer) restErrors.IRestErr
	Count(namespace string) (int, restErrors.IRestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewIpfsPeerService() IService {
	return ipfsPeerService{}
}

// Get gets a single IPFS peer by name
func (service ipfsPeerService) Get(namespacedName types.NamespacedName) (ipfsv1alpha1.Peer, restErrors.IRestErr) {
	peer := &ipfsv1alpha1.Peer{}

	if err := k8sClient.Get(context.Background(), namespacedName, peer); err != nil {
		if apiErrors.IsNotFound(err) {
			return ipfsv1alpha1.Peer{}, restErrors.NewBadRequestError(fmt.Sprintf("peer by name %s doesn't exit", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return ipfsv1alpha1.Peer{}, restErrors.NewInternalServerError(fmt.Sprintf("can't get peer by name %s", peer.Name))
	}

	return *peer, nil
}

// Create creates IPFS peer from spec
func (service ipfsPeerService) Create(dto PeerDto) (ipfsv1alpha1.Peer, restErrors.IRestErr) {
	var initProfiles []ipfsv1alpha1.Profile
	for _, profile := range dto.InitProfiles {
		initProfiles = append(initProfiles, ipfsv1alpha1.Profile(profile))
	}

	peer := &ipfsv1alpha1.Peer{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: ipfsv1alpha1.PeerSpec{
			InitProfiles: initProfiles,
			Image:        dto.Image,
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	k8s.DefaultResources(&peer.Spec.Resources)

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8sClient.Create(context.Background(), peer); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return ipfsv1alpha1.Peer{}, restErrors.NewNotFoundError(fmt.Sprintf("peer by name %s already exits", dto.Name))
		}
		go logger.Error(service.Create, err)
		return ipfsv1alpha1.Peer{}, restErrors.NewInternalServerError("failed to create peer")
	}

	return *peer, nil
}

// Update updates IPFS peer by name from spec
func (service ipfsPeerService) Update(dto PeerDto, peer *ipfsv1alpha1.Peer) restErrors.IRestErr {
	if dto.APIPort != 0 {
		peer.Spec.APIPort = dto.APIPort
	}

	if dto.APIHost != "" {
		peer.Spec.APIHost = dto.APIHost
	}

	if dto.GatewayPort != 0 {
		peer.Spec.GatewayPort = dto.GatewayPort
	}

	if dto.GatewayHost != "" {
		peer.Spec.GatewayHost = dto.GatewayHost
	}

	if dto.Routing != "" {
		peer.Spec.Routing = ipfsv1alpha1.RoutingMechanism(dto.Routing)
	}

	if len(dto.Profiles) != 0 {
		profiles := []ipfsv1alpha1.Profile{}
		for _, profile := range dto.Profiles {
			profiles = append(profiles, ipfsv1alpha1.Profile(profile))
		}
		peer.Spec.Profiles = profiles
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
	if dto.API != nil {
		peer.Spec.API = *dto.API
	}
	if dto.Gateway != nil {
		peer.Spec.Gateway = *dto.Gateway
	}
	if dto.Image != "" {
		peer.Spec.Image = dto.Image
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	pod := &corev1.Pod{}
	podIsPending := false
	if dto.CPU != "" || dto.Memory != "" {
		key := types.NamespacedName{
			Namespace: peer.Namespace,
			Name:      fmt.Sprintf("%s-0", peer.Name),
		}
		err := k8sClient.Get(context.Background(), key, pod)
		if apiErrors.IsNotFound(err) {
			go logger.Error(service.Update, err)
			return restErrors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), peer); err != nil {
		go logger.Error(service.Update, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't update peer by name %s", peer.Name))
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			return restErrors.NewInternalServerError(fmt.Sprintf("can't update peer by name %s", peer.Name))
		}
	}

	return nil
}

// List returns all IPFS peers
func (service ipfsPeerService) List(namespace string) (ipfsv1alpha1.PeerList, restErrors.IRestErr) {
	peers := &ipfsv1alpha1.PeerList{}
	if err := k8sClient.List(context.Background(), peers, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return ipfsv1alpha1.PeerList{}, restErrors.NewInternalServerError("failed to get all peers")
	}

	return *peers, nil
}

// Count returns total number of IPFS peers
func (service ipfsPeerService) Count(namespace string) (int, restErrors.IRestErr) {
	peers := &ipfsv1alpha1.PeerList{}

	if err := k8sClient.List(context.Background(), peers, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.Count, err)
		return 0, restErrors.NewInternalServerError("failed to count all peers")
	}

	return len(peers.Items), nil
}

// Delete deletes ethereum 2.0 IPFS peer by name
func (service ipfsPeerService) Delete(peer *ipfsv1alpha1.Peer) restErrors.IRestErr {
	if err := k8sClient.Delete(context.Background(), peer); err != nil {
		go logger.Error(service.Delete, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delete peer by name %s", peer.Name))
	}

	return nil
}
