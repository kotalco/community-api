package models

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	corev1 "k8s.io/api/core/v1"
)

// Secret is Kubernetes secret
type Secret struct {
	models.Time
	Name string            `json:"name"`
	Type string            `json:"type"`
	Data map[string]string `json:"data,omitempty"`
}

func FromCoreSecret(secret *corev1.Secret) *Secret {
	return &Secret{
		Name: secret.Name,
		Time: models.Time{
			CreatedAt: secret.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Type: secret.Labels["kotal.io/key-type"],
	}
}
