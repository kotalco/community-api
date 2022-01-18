package models

import (
	"github.com/kotalco/api/models"
	"github.com/kotalco/api/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
)

// Node is NEAR node
type Node struct {
	models.Time
	Name                     string  `json:"name"`
	Network                  string  `json:"network"`
	Archive                  bool    `json:"archive"`
	NodePrivateKeySecretName string  `json:"nodePrivateKeySecretName"`
	ValidatorSecretName      string  `json:"validatorSecretName"`
	MinPeers                 uint    `json:"minPeers"`
	P2PPort                  uint    `json:"p2pPort"`
	P2PHost                  string  `json:"p2pHost"`
	RPC                      *bool   `json:"rpc"`
	RPCPort                  uint    `json:"rpcPort"`
	RPCHost                  string  `json:"rpcHost"`
	CPU                      string  `json:"cpu"`
	CPULimit                 string  `json:"cpuLimit"`
	Memory                   string  `json:"memory"`
	MemoryLimit              string  `json:"memoryLimit"`
	Storage                  string  `json:"storage"`
	StorageClass             *string `json:"storageClass"`
}

// FromNEARNode creates node model from NEAR node
func FromNEARNode(node *nearv1alpha1.Node) *Node {
	return &Node{
		Name: node.Name,
		Time: models.Time{
			CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:                  string(node.Spec.Network),
		Archive:                  node.Spec.Archive,
		NodePrivateKeySecretName: node.Spec.NodePrivateKeySecretName,
		ValidatorSecretName:      node.Spec.ValidatorSecretName,
		MinPeers:                 node.Spec.MinPeers,
		P2PPort:                  node.Spec.P2PPort,
		P2PHost:                  node.Spec.P2PHost,
		RPC:                      &node.Spec.RPC,
		RPCPort:                  node.Spec.RPCPort,
		RPCHost:                  node.Spec.RPCHost,
		CPU:                      node.Spec.CPU,
		CPULimit:                 node.Spec.CPULimit,
		Memory:                   node.Spec.Memory,
		MemoryLimit:              node.Spec.MemoryLimit,
		Storage:                  node.Spec.Storage,
		StorageClass:             node.Spec.StorageClass,
	}
}
