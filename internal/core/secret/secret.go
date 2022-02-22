package secret

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	corev1 "k8s.io/api/core/v1"
)

type SecretDto struct {
	models.Time
	Name string            `json:"name"`
	Type string            `json:"type"`
	Data map[string]string `json:"data,omitempty"`
}

type SecretsDto []SecretDto

func (secret SecretDto) FromCoreSecret(s *corev1.Secret) *SecretDto {
	return &SecretDto{
		Name: s.Name,
		Time: models.Time{
			CreatedAt: s.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Type: s.Labels["kotal.io/key-type"],
	}
}

func (secret SecretsDto) FromCoreSecret(secrets []corev1.Secret) SecretsDto {
	result := make(SecretsDto, len(secrets))
	for index, value := range secrets {
		result[index] = *(SecretDto{}.FromCoreSecret(&value))
	}
	return result
}
