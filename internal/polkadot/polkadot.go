package polkadot

import (
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
)

type PolkadotDto struct {
	Name                     string   `json:"name"`
	Network                  string   `json:"network"`
	NodePrivateKeySecretName string   `json:"nodePrivateKeySecretName"`
	Validator                *bool    `json:"validator"`
	SyncMode                 string   `json:"syncMode"`
	P2PPort                  uint     `json:"p2pPort"`
	Pruning                  *bool    `json:"pruning"`
	RetainedBlocks           uint     `json:"retainedBlocks"`
	Logging                  string   `json:"logging"`
	Telemetry                *bool    `json:"telemetry"`
	TelemetryURL             string   `json:"telemetryURL"`
	Prometheus               *bool    `json:"prometheus"`
	PrometheusPort           uint     `json:"prometheusPort"`
	RPC                      *bool    `json:"rpc"`
	RPCPort                  uint     `json:"rpcPort"`
	WS                       *bool    `json:"ws"`
	WSPort                   uint     `json:"wsPort"`
	CORSDomains              []string `json:"corsDomains"`
	CPU                      string   `json:"cpu"`
	CPULimit                 string   `json:"cpuLimit"`
	Memory                   string   `json:"memory"`
	MemoryLimit              string   `json:"memoryLimit"`
	Storage                  string   `json:"storage"`
	StorageClass             *string  `json:"storageClass"`
}

type PolkadotListDto []PolkadotDto

func (dto PolkadotDto) FromPolkadotNode(node *polkadotv1alpha1.Node) *PolkadotDto {
	return &PolkadotDto{
		Name:                     node.Name,
		Network:                  node.Spec.Network,
		NodePrivateKeySecretName: node.Spec.NodePrivateKeySecretName,
		Validator:                &node.Spec.Validator,
		SyncMode:                 string(node.Spec.SyncMode),
		P2PPort:                  node.Spec.P2PPort,
		Pruning:                  node.Spec.Pruning,
		RetainedBlocks:           node.Spec.RetainedBlocks,
		Logging:                  string(node.Spec.Logging),
		Telemetry:                &node.Spec.Telemetry,
		TelemetryURL:             node.Spec.TelemetryURL,
		Prometheus:               &node.Spec.Prometheus,
		PrometheusPort:           node.Spec.PrometheusPort,
		RPC:                      &node.Spec.RPC,
		RPCPort:                  node.Spec.RPCPort,
		WS:                       &node.Spec.WS,
		WSPort:                   node.Spec.WSPort,
		CORSDomains:              node.Spec.CORSDomains,
		CPU:                      node.Spec.CPU,
		CPULimit:                 node.Spec.CPULimit,
		Memory:                   node.Spec.Memory,
		MemoryLimit:              node.Spec.MemoryLimit,
		Storage:                  node.Spec.Storage,
		StorageClass:             node.Spec.StorageClass,
	}
}

func (listDto PolkadotListDto) FromPolkadotNode(nodes []polkadotv1alpha1.Node) PolkadotListDto {
	result := make(PolkadotListDto, len(nodes))
	for index, v := range nodes {
		result[index] = *(PolkadotDto{}.FromPolkadotNode(&v))
	}
	return result
}
