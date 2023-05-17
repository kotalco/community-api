package aptos

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	aptosv1alpha1 "github.com/kotalco/kotal/apis/aptos/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

type AptosDto struct {
	models.Time
	k8s.MetaDataDto
	Network                  aptosv1alpha1.AptosNetwork `json:"network"`
	Image                    string                     `json:"image"`
	Validator                *bool                      `json:"validator"`
	NodePrivateKeySecretName string                     `json:"nodePrivateKeySecretName"`
	API                      *bool                      `json:"api"`
	APIPort                  uint                       `json:"apiPort"`
	P2PPort                  uint                       `json:"p2pPort"`
	sharedAPI.Resources
}

type AptosListDto []AptosDto

func (dto AptosDto) FromAptosNode(n aptosv1alpha1.Node) AptosDto {
	dto.Name = n.Name
	dto.Time = models.Time{CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Image = n.Spec.Image
	dto.Network = n.Spec.Network
	dto.Validator = &n.Spec.Validator
	dto.NodePrivateKeySecretName = n.Spec.NodePrivateKeySecretName
	dto.API = &n.Spec.API
	dto.APIPort = n.Spec.APIPort
	dto.P2PPort = n.Spec.P2PPort
	dto.CPU = n.Spec.CPU
	dto.CPULimit = n.Spec.CPULimit
	dto.Memory = n.Spec.Memory
	dto.MemoryLimit = n.Spec.MemoryLimit
	dto.Storage = n.Spec.Storage
	dto.StorageClass = n.Spec.StorageClass
	return dto
}

func (nodes AptosListDto) FromAptosNode(models []aptosv1alpha1.Node) AptosListDto {
	result := make(AptosListDto, len(models))
	for index, model := range models {
		result[index] = AptosDto{}.FromAptosNode(model)
	}
	return result
}
