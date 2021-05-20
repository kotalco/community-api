package models

import ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"

// Peer is IPFS peer
// TODO: update with SwarmKeySecret and Resources
type Peer struct {
	Name         string   `json:"name"`
	InitProfiles []string `json:"initProfiles"`
	APIPort      uint     `json:"apiPort"`
	APIHost      string   `json:"apiHost"`
	GatewayPort  uint     `json:"gatewayPort"`
	GatewayHost  string   `json:"gatewayHost"`
	Routing      string   `json:"routing"`
	Profiles     []string `json:"profiles"`
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
		Name:         peer.Name,
		APIPort:      peer.Spec.APIPort,
		APIHost:      peer.Spec.APIHost,
		GatewayPort:  peer.Spec.GatewayPort,
		GatewayHost:  peer.Spec.GatewayHost,
		Routing:      string(peer.Spec.Routing),
		Profiles:     profiles,
		InitProfiles: initProfiles,
	}
}
