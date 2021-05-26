package models

import ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"

type Validator struct {
	Name    string `json:"name"`
	Network string `json:"network"`
}

func FromEthereum2Validator(validator *ethereum2v1alpha1.Validator) *Validator {
	return &Validator{
		Name:    validator.Name,
		Network: validator.Spec.Network,
	}
}
