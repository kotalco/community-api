package validator

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type validatorService struct{}

type IService interface {
	Get(types.NamespacedName) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	Create(dto *ValidatorDto) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	Update(*ValidatorDto, *ethereum2v1alpha1.Validator) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	List(namespace string) (*ethereum2v1alpha1.ValidatorList, *errors.RestErr)
	Delete(node *ethereum2v1alpha1.Validator) *errors.RestErr
	Count(namespace string) (*int, *errors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewValidatorService() IService {
	return validatorService{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (service validatorService) Get(namespacedName types.NamespacedName) (*ethereum2v1alpha1.Validator, *errors.RestErr) {

	validator := &ethereum2v1alpha1.Validator{}

	if err := k8sClient.Get(context.Background(), namespacedName, validator); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("validator by name %s doesn't exit", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get a validator by name %s", namespacedName.Name))
	}
	return validator, nil
}

// Create creates ethereum 2.0 beacon node from spec
func (service validatorService) Create(dto *ValidatorDto) (*ethereum2v1alpha1.Validator, *errors.RestErr) {
	keystores := []ethereum2v1alpha1.Keystore{}
	for _, keystore := range dto.Keystores {
		keystores = append(keystores, ethereum2v1alpha1.Keystore{
			SecretName: keystore.SecretName,
		})
	}

	validator := &ethereum2v1alpha1.Validator{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: ethereum2v1alpha1.ValidatorSpec{
			Network:   dto.Network,
			Client:    ethereum2v1alpha1.Ethereum2Client(dto.Client),
			Keystores: keystores,
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if dto.Client == string(ethereum2v1alpha1.PrysmClient) && dto.WalletPasswordSecretName != "" {
		validator.Spec.WalletPasswordSecret = dto.WalletPasswordSecretName
	}

	if len(dto.BeaconEndpoints) != 0 {
		validator.Spec.BeaconEndpoints = dto.BeaconEndpoints
	} else {
		validator.Spec.BeaconEndpoints = []string{}
	}

	if os.Getenv("MOCK") == "true" {
		validator.Default()
	}

	if err := k8sClient.Create(context.Background(), validator); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("validator by name %s already exits", validator.Name))
		}
		go logger.Error(service.Create, err)
		return nil, errors.NewInternalServerError("failed to create validator")
	}

	return validator, nil
}

// Update updates ethereum 2.0 beacon node by name from spec
func (service validatorService) Update(dto *ValidatorDto, validator *ethereum2v1alpha1.Validator) (*ethereum2v1alpha1.Validator, *errors.RestErr) {
	if dto.WalletPasswordSecretName != "" {
		validator.Spec.WalletPasswordSecret = dto.WalletPasswordSecretName
	}

	if len(dto.Keystores) != 0 {
		keystores := []ethereum2v1alpha1.Keystore{}
		for _, keystore := range dto.Keystores {
			keystores = append(keystores, ethereum2v1alpha1.Keystore{
				SecretName: keystore.SecretName,
			})
		}
		validator.Spec.Keystores = keystores
	}

	if dto.Graffiti != "" {
		validator.Spec.Graffiti = dto.Graffiti
	}

	if len(dto.BeaconEndpoints) != 0 {
		validator.Spec.BeaconEndpoints = dto.BeaconEndpoints
	}

	if dto.CPU != "" {
		validator.Spec.CPU = dto.CPU
	}
	if dto.CPULimit != "" {
		validator.Spec.CPULimit = dto.CPULimit
	}
	if dto.Memory != "" {
		validator.Spec.Memory = dto.Memory
	}
	if dto.MemoryLimit != "" {
		validator.Spec.MemoryLimit = dto.MemoryLimit
	}
	if dto.Storage != "" {
		validator.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		validator.Default()
	}

	if err := k8sClient.Update(context.Background(), validator); err != nil {
		go logger.Error(service.Update, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update validator by name %s", validator.Name))
	}

	return validator, nil
}

// List returns all ethereum 2.0 beacon nodes
func (service validatorService) List(namespace string) (*ethereum2v1alpha1.ValidatorList, *errors.RestErr) {
	validators := &ethereum2v1alpha1.ValidatorList{}

	if err := k8sClient.List(context.Background(), validators, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all validators")
	}

	return validators, nil
}

// Count returns total number of beacon nodes
func (service validatorService) Count(namespace string) (*int, *errors.RestErr) {
	validators := &ethereum2v1alpha1.ValidatorList{}

	if err := k8sClient.List(context.Background(), validators, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.Count, err)
		return nil, errors.NewInternalServerError("error counting validators")
	}

	length := len(validators.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 beacon node by name
func (service validatorService) Delete(validator *ethereum2v1alpha1.Validator) *errors.RestErr {
	if err := k8sClient.Delete(context.Background(), validator); err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewBadRequestError(fmt.Sprintf("can't delete validator by name %s", validator.Name))
	}

	return nil
}
