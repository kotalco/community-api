package models

import ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"

type ClusterPeer struct {
	Name                 string   `json:"name"`
	ID                   string   `json:"id"`
	PrivatekeySecretName string   `json:"privatekeySecretName"`
	TrustedPeers         []string `json:"trustedPeers"`
}

func FromIPFSClusterPeer(peer *ipfsv1alpha1.ClusterPeer) *ClusterPeer {
	return &ClusterPeer{
		Name:                 peer.Name,
		ID:                   peer.Spec.ID,
		PrivatekeySecretName: peer.Spec.PrivatekeySecretName,
		TrustedPeers:         peer.Spec.TrustedPeers,
	}
}
