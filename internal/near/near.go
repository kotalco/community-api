package near

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
)

// NearDto is NEAR node
type NearDto struct {
	models.Time
	Name                     string    `json:"name"`
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
	CPU                      string    `json:"cpu"`
	CPULimit                 string    `json:"cpuLimit"`
	Memory                   string    `json:"memory"`
	MemoryLimit              string    `json:"memoryLimit"`
	Storage                  string    `json:"storage"`
	StorageClass             *string   `json:"storageClass"`
}

type NearListDto []NearDto

// FromNEARNode creates node model from NEAR node
func (dto NearDto) FromNEARNode(node *nearv1alpha1.Node) *NearDto {
	return &NearDto{
		Name: node.Name,
		Time: models.Time{
			CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:                  string(node.Spec.Network),
		Archive:                  node.Spec.Archive,
		NodePrivateKeySecretName: node.Spec.NodePrivateKeySecretName,
		ValidatorSecretName:      node.Spec.ValidatorSecretName,
		MinPeers:                 node.Spec.MinPeers,
		P2PPort:                  node.Spec.P2PPort,
		P2PHost:                  node.Spec.P2PHost,
		RPC:                      &node.Spec.RPC,
		RPCPort:                  node.Spec.RPCPort,
		RPCHost:                  node.Spec.RPCHost,
		PrometheusPort:           node.Spec.PrometheusPort,
		PrometheusHost:           node.Spec.PrometheusHost,
		TelemetryURL:             node.Spec.TelemetryURL,
		Bootnodes:                &node.Spec.Bootnodes,
		CPU:                      node.Spec.CPU,
		CPULimit:                 node.Spec.CPULimit,
		Memory:                   node.Spec.Memory,
		MemoryLimit:              node.Spec.MemoryLimit,
		Storage:                  node.Spec.Storage,
		StorageClass:             node.Spec.StorageClass,
	}
}

func (listDto NearListDto) FromNEARNode(nodes []nearv1alpha1.Node) NearListDto {
	result := make(NearListDto, len(nodes))
	for index, v := range nodes {
		result[index] = *(NearDto{}.FromNEARNode(&v))
	}
	return result
}
