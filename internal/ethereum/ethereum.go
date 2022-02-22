package ethereum

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
)

// ImportedAccount is account derived from private key
type ImportedAccount struct {
	PrivateKeySecretName string `json:"privateKeySecretName"`
	PasswordSecretName   string `json:"passwordSecretName"`
}

// Node is Ethereum node
type EthereumDto struct {
	models.Time
	Name                     string           `json:"name"`
	Network                  string           `json:"network"`
	Client                   string           `json:"client"`
	Logging                  string           `json:"logging"`
	NodePrivateKeySecretName string           `json:"nodePrivateKeySecretName"`
	SyncMode                 string           `json:"syncMode"`
	P2PPort                  uint             `json:"p2pPort"`
	StaticNodes              *[]string        `json:"staticNodes"`
	Bootnodes                *[]string        `json:"bootnodes"`
	Miner                    *bool            `json:"miner"`
	Coinbase                 string           `json:"coinbase"`
	Import                   *ImportedAccount `json:"import"`
	RPC                      *bool            `json:"rpc"`
	RPCPort                  uint             `json:"rpcPort"`
	RPCAPI                   []string         `json:"rpcAPI"`
	WS                       *bool            `json:"ws"`
	WSPort                   uint             `json:"wsPort"`
	WSAPI                    []string         `json:"wsAPI"`
	GraphQL                  *bool            `json:"graphql"`
	GraphQLPort              uint             `json:"graphqlPort"`
	Hosts                    []string         `json:"hosts"`
	CORSDomains              []string         `json:"corsDomains"`
	CPU                      string           `json:"cpu"`
	CPULimit                 string           `json:"cpuLimit"`
	Memory                   string           `json:"memory"`
	MemoryLimit              string           `json:"memoryLimit"`
	Storage                  string           `json:"storage"`
	StorageClass             *string          `json:"storageClass"`
}
type EthereumListDto []EthereumDto

func (node EthereumDto) FromEthereumNode(n *ethereumv1alpha1.Node) *EthereumDto {
	model := &EthereumDto{
		Name: n.Name,
		Time: models.Time{
			CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:                  n.Spec.Network,
		Client:                   string(n.Spec.Client),
		Logging:                  string(n.Spec.Logging),
		NodePrivateKeySecretName: n.Spec.NodePrivateKeySecretName,
		SyncMode:                 string(n.Spec.SyncMode),
		P2PPort:                  n.Spec.P2PPort,
		Miner:                    &n.Spec.Miner,
		Coinbase:                 string(n.Spec.Coinbase),
		RPC:                      &n.Spec.RPC,
		RPCPort:                  n.Spec.RPCPort,
		WS:                       &n.Spec.WS,
		WSPort:                   n.Spec.WSPort,
		GraphQL:                  &n.Spec.GraphQL,
		GraphQLPort:              n.Spec.GraphQLPort,
		Hosts:                    n.Spec.Hosts,
		CORSDomains:              n.Spec.CORSDomains,
		CPU:                      n.Spec.CPU,
		CPULimit:                 n.Spec.CPULimit,
		Memory:                   n.Spec.Memory,
		MemoryLimit:              n.Spec.MemoryLimit,
		Storage:                  n.Spec.Storage,
		StorageClass:             n.Spec.StorageClass,
	}

	if n.Spec.Miner && n.Spec.Import != nil {
		model.Import = &ImportedAccount{
			PrivateKeySecretName: n.Spec.Import.PrivateKeySecretName,
			PasswordSecretName:   n.Spec.Import.PasswordSecretName,
		}
	}

	var rpcAPI []string
	for _, api := range n.Spec.RPCAPI {
		rpcAPI = append(rpcAPI, string(api))
	}
	model.RPCAPI = rpcAPI

	var wsAPI []string
	for _, api := range n.Spec.WSAPI {
		wsAPI = append(wsAPI, string(api))
	}
	model.WSAPI = wsAPI

	staticNodes := []string{}
	for _, enode := range n.Spec.StaticNodes {
		staticNodes = append(staticNodes, string(enode))
	}
	model.StaticNodes = &staticNodes

	bootnodes := []string{}
	for _, bootnode := range n.Spec.Bootnodes {
		bootnodes = append(bootnodes, string(bootnode))
	}
	model.Bootnodes = &bootnodes

	return model
}

func (nodes EthereumListDto) FromEthereumNode(models []ethereumv1alpha1.Node) EthereumListDto {
	result := make(EthereumListDto, len(models))
	for index, v := range models {
		result[index] = *(EthereumDto{}.FromEthereumNode(&v))
	}
	return result
}
