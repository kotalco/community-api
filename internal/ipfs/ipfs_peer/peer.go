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
	Name         string   `json:"name"`
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
func (peerDto PeerDto) FromIPFSPeer(peer *ipfsv1alpha1.Peer) *PeerDto {
	var profiles, initProfiles []string

	// init profiles
	for _, profile := range peer.Spec.InitProfiles {
		initProfiles = append(initProfiles, string(profile))
	}

	// profiles
	for _, profile := range peer.Spec.Profiles {
		profiles = append(profiles, string(profile))
	}

	return &PeerDto{
		Name: peer.Name,
		Time: models.Time{
			CreatedAt: peer.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		APIPort:      peer.Spec.APIPort,
		APIHost:      peer.Spec.APIHost,
		GatewayPort:  peer.Spec.GatewayPort,
		GatewayHost:  peer.Spec.GatewayHost,
		Routing:      string(peer.Spec.Routing),
		Profiles:     profiles,
		InitProfiles: initProfiles,
		CPU:          peer.Spec.CPU,
		CPULimit:     peer.Spec.CPULimit,
		Memory:       peer.Spec.Memory,
		MemoryLimit:  peer.Spec.MemoryLimit,
		Storage:      peer.Spec.Storage,
		StorageClass: peer.Spec.StorageClass,
	}
}

func (peerListDto PeerListDto) FromIPFSPeer(peers []ipfsv1alpha1.Peer) PeerListDto {
	result := make(PeerListDto, len(peers))
	for index, v := range peers {
		result[index] = *(PeerDto{}.FromIPFSPeer(&v))
	}
	return result
}
