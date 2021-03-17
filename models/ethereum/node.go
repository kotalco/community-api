package params

import ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"

// Node is Ethereum node
type Node struct {
	Name    string `json:"name"`
	Network string `json:"network"`
	Client  string `json:"client"`
}

func FromEthereumNode(n *ethereumv1alpha1.Node) *Node {
	return &Node{
		Name:    n.Name,
		Network: n.Spec.Join,
		Client:  string(n.Spec.Client),
	}
}
