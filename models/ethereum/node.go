package params

import ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"

// Node is Ethereum node
type Node struct {
	Name    string   `json:"name"`
	Network string   `json:"network"`
	Client  string   `json:"client"`
	RPC     bool     `json:"rpc"`
	RPCAPI  []string `json:"rpcAPI"`
	WS      bool     `json:"ws"`
}

func FromEthereumNode(n *ethereumv1alpha1.Node) *Node {
	model := &Node{
		Name:    n.Name,
		Network: n.Spec.Join,
		Client:  string(n.Spec.Client),
		RPC:     n.Spec.RPC,
		WS:      n.Spec.WS,
	}

	// convert ethereum API into string
	rpcAPI := []string{}
	for _, api := range n.Spec.RPCAPI {
		rpcAPI = append(rpcAPI, string(api))
	}

	model.RPCAPI = rpcAPI

	return model

}
