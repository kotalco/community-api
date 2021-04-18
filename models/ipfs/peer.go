package models

import ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"

type Peer struct {
	Name        string `json:"name"`
	APIPort     uint   `json:"apiPort"`
	APIHost     string `json:"apiHost"`
	GatewayPort uint   `json:"gatewayPort"`
	GatewayHost string `json:"gatewayHost"`
}

// FromIPFSPeer creates peer model from IPFS peer
func FromIPFSPeer(peer *ipfsv1alpha1.Peer) *Peer {
	return &Peer{
		Name:        peer.Name,
		APIPort:     peer.Spec.APIPort,
		APIHost:     peer.Spec.APIHost,
		GatewayPort: peer.Spec.GatewayPort,
		GatewayHost: peer.Spec.GatewayHost,
	}
}
