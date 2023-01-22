package ethereum

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

// ImportedAccount is account derived from private key
type ImportedAccount struct {
	PrivateKeySecretName string `json:"privateKeySecretName"`
	PasswordSecretName   string `json:"passwordSecretName"`
}

// Node is Ethereum node
type EthereumDto struct {
	models.Time
	k8s.MetaDataDto
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
	Engine                   *bool            `json:"engine"`
	sharedAPI.Resources
}
type EthereumListDto []EthereumDto

func (dto EthereumDto) FromEthereumNode(node *ethereumv1alpha1.Node) *EthereumDto {
	dto.Name = node.Name
	dto.Time = models.Time{CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = node.Spec.Network
	dto.Client = string(node.Spec.Client)
	dto.Logging = string(node.Spec.Logging)
	dto.NodePrivateKeySecretName = node.Spec.NodePrivateKeySecretName
	dto.SyncMode = string(node.Spec.SyncMode)
	dto.P2PPort = node.Spec.P2PPort
	dto.Miner = &node.Spec.Miner
	dto.Coinbase = string(node.Spec.Coinbase)
	dto.RPC = &node.Spec.RPC
	dto.RPCPort = node.Spec.RPCPort
	dto.WS = &node.Spec.WS
	dto.WSPort = node.Spec.WSPort
	dto.GraphQL = &node.Spec.GraphQL
	dto.GraphQLPort = node.Spec.GraphQLPort
	dto.Hosts = node.Spec.Hosts
	dto.CORSDomains = node.Spec.CORSDomains
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass
	dto.Engine = &node.Spec.Engine

	if node.Spec.Miner && node.Spec.Import != nil {
		dto.Import = &ImportedAccount{
			PrivateKeySecretName: node.Spec.Import.PrivateKeySecretName,
			PasswordSecretName:   node.Spec.Import.PasswordSecretName,
		}
	}

	var rpcAPI []string
	for _, api := range node.Spec.RPCAPI {
		rpcAPI = append(rpcAPI, string(api))
	}
	dto.RPCAPI = rpcAPI

	var wsAPI []string
	for _, api := range node.Spec.WSAPI {
		wsAPI = append(wsAPI, string(api))
	}
	dto.WSAPI = wsAPI

	staticNodes := []string{}
	for _, enode := range node.Spec.StaticNodes {
		staticNodes = append(staticNodes, string(enode))
	}
	dto.StaticNodes = &staticNodes

	bootnodes := []string{}
	for _, bootnode := range node.Spec.Bootnodes {
		bootnodes = append(bootnodes, string(bootnode))
	}
	dto.Bootnodes = &bootnodes

	return &dto
}

func (nodes EthereumListDto) FromEthereumNode(models []ethereumv1alpha1.Node) EthereumListDto {
	result := make(EthereumListDto, len(models))
	for index, v := range models {
		result[index] = *(EthereumDto{}.FromEthereumNode(&v))
	}
	return result
}
