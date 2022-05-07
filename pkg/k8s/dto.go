package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type MetaDataDto struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func (workspace *MetaDataDto) ObjectMetaFromNamespaceDto() metav1.ObjectMeta {
	if workspace.Namespace == "" {
		workspace.Namespace = "default"
	}
	return metav1.ObjectMeta{
		Name:      workspace.Name,
		Namespace: workspace.Namespace,
	}
}
