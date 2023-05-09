package near

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

// NearDto is NEAR node
type NearDto struct {
	models.Time
	k8s.MetaDataDto
	Network                  string    `json:"network"`
	Archive                  bool      `json:"archive"`
	NodePrivateKeySecretName string    `json:"nodePrivateKeySecretName"`
	ValidatorSecretName      string    `json:"validatorSecretName"`
	MinPeers                 uint      `json:"minPeers"`
	P2PPort                  uint      `json:"p2pPort"`
	P2PHost                  string    `json:"p2pHost"`
	RPC                      *bool     `json:"rpc"`
	RPCPort                  uint      `json:"rpcPort"`
	RPCHost                  string    `json:"rpcHost"`
	PrometheusPort           uint      `json:"prometheusPort"`
	PrometheusHost           string    `json:"prometheusHost"`
	TelemetryURL             string    `json:"telemetryURL"`
	Bootnodes                *[]string `json:"bootnodes"`
	Image                    string    `json:"image"`
	sharedAPI.Resources
}

type NearListDto []NearDto

// FromNEARNode creates node model from NEAR node
func (dto NearDto) FromNEARNode(node nearv1alpha1.Node) NearDto {
	dto.Name = node.Name
	dto.Time = models.Time{CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = string(node.Spec.Network)
	dto.Archive = node.Spec.Archive
	dto.NodePrivateKeySecretName = node.Spec.NodePrivateKeySecretName
	dto.ValidatorSecretName = node.Spec.ValidatorSecretName
	dto.MinPeers = node.Spec.MinPeers
	dto.P2PPort = node.Spec.P2PPort
	dto.P2PHost = node.Spec.P2PHost
	dto.RPC = &node.Spec.RPC
	dto.RPCPort = node.Spec.RPCPort
	dto.RPCHost = node.Spec.RPCHost
	dto.PrometheusPort = node.Spec.PrometheusPort
	dto.PrometheusHost = node.Spec.PrometheusHost
	dto.TelemetryURL = node.Spec.TelemetryURL
	dto.Bootnodes = &node.Spec.Bootnodes
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass
	dto.Image = node.Spec.Image

	return dto
}

func (listDto NearListDto) FromNEARNode(nodes []nearv1alpha1.Node) NearListDto {
	result := make(NearListDto, len(nodes))
	for index, v := range nodes {
		result[index] = NearDto{}.FromNEARNode(v)
	}
	return result
}
