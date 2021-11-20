package models

import chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"

type Node struct {
	Name string `json:"name"`
}

func FromChainlinkNode(node *chainlinkv1alpha1.Node) *Node {
	return &Node{
		Name: node.Name,
	}
}
