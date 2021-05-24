package models

import ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"

type BeaconNode struct {
	Name    string `json:"name"`
	Network string `json:"network"`
}

func FromEthereum2BeaconNode(beaconnode *ethereum2v1alpha1.BeaconNode) *BeaconNode {
	return &BeaconNode{
		Name:    beaconnode.Name,
		Network: beaconnode.Spec.Join,
	}
}
