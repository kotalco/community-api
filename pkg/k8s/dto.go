package k8s

import (
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetaDataDto struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func (metaDto *MetaDataDto) ObjectMetaFromMetadataDto() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      metaDto.Name,
		Namespace: metaDto.Namespace,
	}
}

func DefaultResources(res *sharedAPI.Resources) {
	res.CPU = "1"
	res.Memory = "1Gi"
}
