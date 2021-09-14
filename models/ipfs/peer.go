package models

import (
	"github.com/kotalco/api/models"
	"github.com/kotalco/api/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
)

// Peer is IPFS peer
// TODO: update with SwarmKeySecret and Resources
type Peer struct {
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
}

// FromIPFSPeer creates peer model from IPFS peer
func FromIPFSPeer(peer *ipfsv1alpha1.Peer) *Peer {
	var profiles, initProfiles []string

	// init profiles
	for _, profile := range peer.Spec.InitProfiles {
		initProfiles = append(initProfiles, string(profile))
	}

	// profiles
	for _, profile := range peer.Spec.Profiles {
		profiles = append(profiles, string(profile))
	}

	return &Peer{
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
	}
}
