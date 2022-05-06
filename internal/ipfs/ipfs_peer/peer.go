package ipfs_peer

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
)

// Peer is IPFS peer
// TODO: update with SwarmKeySecret and Resources
type PeerDto struct {
	models.Time
	models.NamespaceDto
	InitProfiles []string `json:"initProfiles"`
	APIPort      uint     `json:"apiPort"`
	APIHost      string   `json:"apiHost"`
	GatewayPort  uint     `json:"gatewayPort"`
	GatewayHost  string   `json:"gatewayHost"`
	Routing      string   `json:"routing"`
	Profiles     []string `json:"profiles"`
	CPU          string   `json:"cpu"`
	CPULimit     string   `json:"cpuLimit"`
	Memory       string   `json:"memory"`
	MemoryLimit  string   `json:"memoryLimit"`
	Storage      string   `json:"storage"`
	StorageClass *string  `json:"storageClass"`
}

type PeerListDto []PeerDto

// FromIPFSPeer creates peer model from IPFS peer
func (dto PeerDto) FromIPFSPeer(peer *ipfsv1alpha1.Peer) *PeerDto {
	var profiles, initProfiles []string

	// init profiles
	for _, profile := range peer.Spec.InitProfiles {
		initProfiles = append(initProfiles, string(profile))
	}

	// profiles
	for _, profile := range peer.Spec.Profiles {
		profiles = append(profiles, string(profile))
	}

	dto.Name = peer.Name
	dto.Time = models.Time{CreatedAt: peer.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.APIPort = peer.Spec.APIPort
	dto.APIHost = peer.Spec.APIHost
	dto.GatewayPort = peer.Spec.GatewayPort
	dto.GatewayHost = peer.Spec.GatewayHost
	dto.Routing = string(peer.Spec.Routing)
	dto.Profiles = profiles
	dto.InitProfiles = initProfiles
	dto.CPU = peer.Spec.CPU
	dto.CPULimit = peer.Spec.CPULimit
	dto.Memory = peer.Spec.Memory
	dto.MemoryLimit = peer.Spec.MemoryLimit
	dto.Storage = peer.Spec.Storage
	dto.StorageClass = peer.Spec.StorageClass

	return &dto
}

func (peerListDto PeerListDto) FromIPFSPeer(peers []ipfsv1alpha1.Peer) PeerListDto {
	result := make(PeerListDto, len(peers))
	for index, v := range peers {
		result[index] = *(PeerDto{}.FromIPFSPeer(&v))
	}
	return result
}
