package validator

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

type ValidatorDto struct {
	models.Time
	Name                     string        `json:"name"`
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

func (validatorDto ValidatorDto) FromEthereum2Validator(validator *ethereum2v1alpha1.Validator) *ValidatorDto {
	keystores := []KeystoreDto{}
	for _, keystore := range validator.Spec.Keystores {
		keystores = append(keystores, KeystoreDto{
			SecretName: keystore.SecretName,
		})
	}
	return &ValidatorDto{
		Name: validator.Name,
		Time: models.Time{
			CreatedAt: validator.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:                  validator.Spec.Network,
		Client:                   string(validator.Spec.Client),
		Graffiti:                 validator.Spec.Graffiti,
		BeaconEndpoints:          validator.Spec.BeaconEndpoints,
		Keystores:                keystores,
		CPU:                      validator.Spec.CPU,
		CPULimit:                 validator.Spec.CPULimit,
		Memory:                   validator.Spec.Memory,
		MemoryLimit:              validator.Spec.MemoryLimit,
		Storage:                  validator.Spec.Storage,
		StorageClass:             validator.Spec.StorageClass,
		WalletPasswordSecretName: validator.Spec.WalletPasswordSecret,
	}
}

func (validatorListDto ValidatorListDto) FromEthereum2Validator(validators []ethereum2v1alpha1.Validator) ValidatorListDto {
	result := make(ValidatorListDto, len(validators))
	for index, v := range validators {
		result[index] = *(ValidatorDto{}.FromEthereum2Validator(&v))
	}
	return result
}
