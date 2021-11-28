package models

import polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"

type Node struct {
	Name                     string `json:"name"`
	Network                  string `json:"network"`
	NodePrivateKeySecretName string `json:"nodePrivateKeySecretName"`
	Validator                *bool  `json:"validator"`
	SyncMode                 string `json:"syncMode"`
	Pruning                  *bool  `json:"pruning"`
	RetainedBlocks           uint   `json:"retainedBlocks"`
	Logging                  string `json:"logging"`
}

func FromPolkadotNode(node *polkadotv1alpha1.Node) *Node {
	return &Node{
		Name:                     node.Name,
		Network:                  node.Spec.Network,
		NodePrivateKeySecretName: node.Spec.NodePrivateKeySecretName,
		Validator:                &node.Spec.Validator,
		SyncMode:                 string(node.Spec.SyncMode),
		Pruning:                  node.Spec.Pruning,
		RetainedBlocks:           node.Spec.RetainedBlocks,
		Logging:                  string(node.Spec.Logging),
	}
}
