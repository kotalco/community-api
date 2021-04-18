package models

import ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"

type Peer struct {
	Name string `json:"name"`
}

// FromIPFSPeer creates peer model from IPFS peer
func FromIPFSPeer(peer *ipfsv1alpha1.Peer) *Peer {
	return &Peer{
		Name: peer.Name,
	}
}
