package ipfs_cluster_peer

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

type ClusterPeerDto struct {
	models.Time
	k8s.MetaDataDto
	ID                   string   `json:"id"`
	PrivatekeySecretName string   `json:"privatekeySecretName"`
	TrustedPeers         []string `json:"trustedPeers"`
	BootstrapPeers       []string `json:"bootstrapPeers"`
	Consensus            string   `json:"consensus"`
	ClusterSecretName    *string  `json:"clusterSecretName"`
	PeerEndpoint         string   `json:"peerEndpoint"`
	Image                string   `json:"image"`
	sharedAPI.Resources
}
type ClusterPeerListDto []ClusterPeerDto

func (dto ClusterPeerDto) FromIPFSClusterPeer(peer ipfsv1alpha1.ClusterPeer) ClusterPeerDto {
	dto.Name = peer.Name
	dto.Time = models.Time{CreatedAt: peer.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.ID = peer.Spec.ID
	dto.PrivatekeySecretName = peer.Spec.PrivateKeySecretName
	dto.TrustedPeers = peer.Spec.TrustedPeers
	dto.BootstrapPeers = peer.Spec.BootstrapPeers
	dto.Consensus = string(peer.Spec.Consensus)
	dto.ClusterSecretName = &peer.Spec.ClusterSecretName
	dto.PeerEndpoint = peer.Spec.PeerEndpoint
	dto.CPU = peer.Spec.CPU
	dto.CPULimit = peer.Spec.CPULimit
	dto.Memory = peer.Spec.Memory
	dto.MemoryLimit = peer.Spec.MemoryLimit
	dto.Storage = peer.Spec.Storage
	dto.StorageClass = peer.Spec.StorageClass
	dto.Image = peer.Spec.Image

	return dto
}

func (clusterPeerListDto ClusterPeerListDto) FromIPFSClusterPeer(peers []ipfsv1alpha1.ClusterPeer) ClusterPeerListDto {
	result := make(ClusterPeerListDto, len(peers))
	for index, v := range peers {
		result[index] = ClusterPeerDto{}.FromIPFSClusterPeer(v)
	}
	return result
}
