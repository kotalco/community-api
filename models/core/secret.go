package models

import corev1 "k8s.io/api/core/v1"

// Secret is Kubernetes secret
type Secret struct {
	Name string            `json:"name"`
	Type string            `json:"type"`
	Data map[string]string `json:"data,omitempty"`
}

func FromCoreSecret(secret *corev1.Secret) *Secret {
	return &Secret{
		Name: secret.Name,
		Type: secret.Labels["kotal.io/key-type"],
	}
}
