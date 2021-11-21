package models

import chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"

type Node struct {
	Name                string `json:"name"`
	EthereumChainId     uint   `json:"ethereumChainId"`
	LinkContractAddress string `json:"linkContractAddress"`
}

func FromChainlinkNode(node *chainlinkv1alpha1.Node) *Node {
	return &Node{
		Name:                node.Name,
		EthereumChainId:     node.Spec.EthereumChainId,
		LinkContractAddress: node.Spec.LinkContractAddress,
	}
}
