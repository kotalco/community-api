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
	dto.Name = node.Name
	dto.Network = node.Spec.Network
	dto.NodePrivateKeySecretName = node.Spec.NodePrivateKeySecretName
	dto.Validator = &node.Spec.Validator
	dto.SyncMode = string(node.Spec.SyncMode)
	dto.P2PPort = node.Spec.P2PPort
	dto.Pruning = node.Spec.Pruning
	dto.RetainedBlocks = node.Spec.RetainedBlocks
	dto.Logging = string(node.Spec.Logging)
	dto.Telemetry = &node.Spec.Telemetry
	dto.TelemetryURL = node.Spec.TelemetryURL
	dto.Prometheus = &node.Spec.Prometheus
	dto.PrometheusPort = node.Spec.PrometheusPort
	dto.RPC = &node.Spec.RPC
	dto.RPCPort = node.Spec.RPCPort
	dto.WS = &node.Spec.WS
	dto.WSPort = node.Spec.WSPort
	dto.CORSDomains = node.Spec.CORSDomains
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass

	return &dto
}

func (listDto PolkadotListDto) FromPolkadotNode(nodes []polkadotv1alpha1.Node) PolkadotListDto {
	result := make(PolkadotListDto, len(nodes))
	for index, v := range nodes {
		result[index] = *(PolkadotDto{}.FromPolkadotNode(&v))
	}
	return result
}
