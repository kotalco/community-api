package stacks

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	stacksv1alpha1 "github.com/kotalco/kotal/apis/stacks/v1alpha1"
)

type StacksDto struct {
	models.Time
	k8s.MetaDataDto
	Image   string                       `json:"image"`
	Network stacksv1alpha1.StacksNetwork `json:"network"`
	P2PPort uint                         `json:"p2pPort"`
	RPCPort uint                         `json:"rpcPort"`
	sharedAPI.Resources
}

type StacksListDto []StacksDto

func (dto StacksDto) FromStacksNode(n *stacksv1alpha1.Node) *StacksDto {
	dto.Name = n.Name
	dto.Time = models.Time{CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Image = n.Spec.Image
	dto.Network = n.Spec.Network
	dto.P2PPort = n.Spec.P2PPort
	dto.RPCPort = n.Spec.RPCPort
	dto.CPU = n.Spec.CPU
	dto.CPULimit = n.Spec.CPULimit
	dto.Memory = n.Spec.Memory
	dto.MemoryLimit = n.Spec.MemoryLimit
	dto.Storage = n.Spec.Storage
	dto.StorageClass = n.Spec.StorageClass
	return &dto
}

func (nodes StacksListDto) FromStacksNode(models []stacksv1alpha1.Node) StacksListDto {
	result := make(StacksListDto, len(models))
	for index, model := range models {
		result[index] = *(StacksDto{}.FromStacksNode(&model))
	}
	return result
}
