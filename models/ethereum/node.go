package models

import (
	"github.com/kotalco/api/models"
	"github.com/kotalco/api/shared"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
)

// Node is Ethereum node
type Node struct {
	models.Time
	Name        string   `json:"name"`
	Network     string   `json:"network"`
	Client      string   `json:"client"`
	SyncMode    string   `json:"syncMode"`
	P2PPort     uint     `json:"p2pPort"`
	RPC         *bool    `json:"rpc"`
	RPCPort     uint     `json:"rpcPort"`
	RPCAPI      []string `json:"rpcAPI"`
	WS          *bool    `json:"ws"`
	WSPort      uint     `json:"wsPort"`
	WSAPI       []string `json:"wsAPI"`
	GraphQL     *bool    `json:"graphql"`
	GraphQLPort uint     `json:"graphqlPort"`
	Hosts       []string `json:"hosts"`
	CORSDomains []string `json:"corsDomains"`
	CPU         string   `json:"cpu"`
	CPULimit    string   `json:"cpuLimit"`
	Memory      string   `json:"memory"`
	MemoryLimit string   `json:"memoryLimit"`
	Storage     string   `json:"storage"`
}

func FromEthereumNode(n *ethereumv1alpha1.Node) *Node {
	model := &Node{
		Name: n.Name,
		Time: models.Time{
			CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:     n.Spec.Join,
		Client:      string(n.Spec.Client),
		SyncMode:    string(n.Spec.SyncMode),
		P2PPort:     n.Spec.P2PPort,
		RPC:         &n.Spec.RPC,
		WS:          &n.Spec.WS,
		GraphQL:     &n.Spec.GraphQL,
		CPU:         n.Spec.CPU,
		CPULimit:    n.Spec.CPULimit,
		Memory:      n.Spec.Memory,
		MemoryLimit: n.Spec.MemoryLimit,
		Storage:     n.Spec.Storage,
	}

	var rpcAPI []string
	if n.Spec.RPC {
		rpcAPI = []string{}
		for _, api := range n.Spec.RPCAPI {
			rpcAPI = append(rpcAPI, string(api))
		}
		model.RPCPort = n.Spec.RPCPort
		model.RPCAPI = rpcAPI
	}

	var wsAPI []string
	if n.Spec.WS {
		wsAPI = []string{}
		for _, api := range n.Spec.WSAPI {
			wsAPI = append(wsAPI, string(api))
		}
		model.WSPort = n.Spec.WSPort
		model.WSAPI = wsAPI
	}

	if n.Spec.GraphQL {
		model.GraphQLPort = n.Spec.GraphQLPort
	}

	if len(n.Spec.Hosts) != 0 {
		model.Hosts = n.Spec.Hosts
	}

	if len(n.Spec.CORSDomains) != 0 {
		model.CORSDomains = n.Spec.CORSDomains
	}

	return model

}
