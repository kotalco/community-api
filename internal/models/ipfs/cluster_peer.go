package models

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
)

type ClusterPeer struct {
	models.Time
	Name                 string   `json:"name"`
	ID                   string   `json:"id"`
	PrivatekeySecretName string   `json:"privatekeySecretName"`
	TrustedPeers         []string `json:"trustedPeers"`
	BootstrapPeers       []string `json:"bootstrapPeers"`
	Consensus            string   `json:"consensus"`
	ClusterSecretName    string   `json:"clusterSecretName"`
	PeerEndpoint         string   `json:"peerEndpoint"`
	CPU                  string   `json:"cpu"`
	CPULimit             string   `json:"cpuLimit"`
	Memory               string   `json:"memory"`
	MemoryLimit          string   `json:"memoryLimit"`
	Storage              string   `json:"storage"`
	StorageClass         *string  `json:"storageClass"`
}

func FromIPFSClusterPeer(peer *ipfsv1alpha1.ClusterPeer) *ClusterPeer {
	return &ClusterPeer{
		Name: peer.Name,
		Time: models.Time{
			CreatedAt: peer.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		ID:                   peer.Spec.ID,
		PrivatekeySecretName: peer.Spec.PrivateKeySecretName,
		TrustedPeers:         peer.Spec.TrustedPeers,
		BootstrapPeers:       peer.Spec.BootstrapPeers,
		Consensus:            string(peer.Spec.Consensus),
		ClusterSecretName:    peer.Spec.ClusterSecretName,
		PeerEndpoint:         peer.Spec.PeerEndpoint,
		CPU:                  peer.Spec.CPU,
		CPULimit:             peer.Spec.CPULimit,
		Memory:               peer.Spec.Memory,
		MemoryLimit:          peer.Spec.MemoryLimit,
		Storage:              peer.Spec.Storage,
		StorageClass:         peer.Spec.StorageClass,
	}
}
