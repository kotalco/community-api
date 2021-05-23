package models

import ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"

type BeaconNode struct {
	Name string `json:"name"`
}

func FromEthereum2BeaconNode(beaconnode *ethereum2v1alpha1.BeaconNode) *BeaconNode {
	return &BeaconNode{
		Name: beaconnode.Name,
	}
}
