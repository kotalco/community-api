package validator

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

type ValidatorDto struct {
	models.Time
	models.NamespaceDto
	Network                  string        `json:"network"`
	Client                   string        `json:"client"`
	Graffiti                 string        `json:"graffiti"`
	BeaconEndpoints          []string      `json:"beaconEndpoints"`
	WalletPasswordSecretName string        `json:"walletPasswordSecretName"`
	Keystores                []KeystoreDto `json:"keystores"`
	CPU                      string        `json:"cpu"`
	CPULimit                 string        `json:"cpuLimit"`
	Memory                   string        `json:"memory"`
	MemoryLimit              string        `json:"memoryLimit"`
	Storage                  string        `json:"storage"`
	StorageClass             *string       `json:"storageClass"`
}

type KeystoreDto struct {
	SecretName string `json:"secretName"`
}

type ValidatorListDto []ValidatorDto

func (dto ValidatorDto) FromEthereum2Validator(validator *ethereum2v1alpha1.Validator) *ValidatorDto {
	keystores := []KeystoreDto{}
	for _, keystore := range validator.Spec.Keystores {
		keystores = append(keystores, KeystoreDto{
			SecretName: keystore.SecretName,
		})
	}

	dto.Name = validator.Name
	dto.Time = models.Time{CreatedAt: validator.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = validator.Spec.Network
	dto.Client = string(validator.Spec.Client)
	dto.Graffiti = validator.Spec.Graffiti
	dto.BeaconEndpoints = validator.Spec.BeaconEndpoints
	dto.Keystores = keystores
	dto.CPU = validator.Spec.CPU
	dto.CPULimit = validator.Spec.CPULimit
	dto.Memory = validator.Spec.Memory
	dto.MemoryLimit = validator.Spec.MemoryLimit
	dto.Storage = validator.Spec.Storage
	dto.StorageClass = validator.Spec.StorageClass
	dto.WalletPasswordSecretName = validator.Spec.WalletPasswordSecret

	return &dto
}

func (validatorListDto ValidatorListDto) FromEthereum2Validator(validators []ethereum2v1alpha1.Validator) ValidatorListDto {
	result := make(ValidatorListDto, len(validators))
	for index, v := range validators {
		result[index] = *(ValidatorDto{}.FromEthereum2Validator(&v))
	}
	return result
}
