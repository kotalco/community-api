package models

import (
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

type Validator struct {
	Name                     string     `json:"name"`
	Network                  string     `json:"network"`
	Client                   string     `json:"client"`
	Graffiti                 string     `json:"graffiti"`
	BeaconEndpoints          []string   `json:"beaconEndpoints"`
	WalletPasswordSecretName string     `json:"walletPasswordSecretName"`
	Keystores                []Keystore `json:"keystores"`
	CPU                      string     `json:"cpu"`
	CPULimit                 string     `json:"cpuLimit"`
	Memory                   string     `json:"memory"`
	MemoryLimit              string     `json:"memoryLimit"`
	Storage                  string     `json:"storage"`
}

type Keystore struct {
	SecretName string `json:"secretName"`
}

func FromEthereum2Validator(validator *ethereum2v1alpha1.Validator) *Validator {
	keystores := []Keystore{}
	for _, keystore := range validator.Spec.Keystores {
		keystores = append(keystores, Keystore{
			SecretName: keystore.SecretName,
		})
	}

	return &Validator{
		Name:                     validator.Name,
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
		WalletPasswordSecretName: validator.Spec.WalletPasswordSecret,
	}
}
