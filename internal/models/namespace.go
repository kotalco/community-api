package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceDto struct {
	Name     string     `json:"name"`
	Namespace string    `json:"namespace,omitempty"`
}

func (workspace *NamespaceDto)ObjectMetaFromNamespaceDto() metav1.ObjectMeta {
	if workspace.Namespace=="" {
		workspace.Namespace="default"
	}
	return metav1.ObjectMeta{
		Name:      workspace.Name,
		Namespace: workspace.Namespace,
	}
}


