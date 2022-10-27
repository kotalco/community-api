package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
