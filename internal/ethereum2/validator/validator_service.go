package validator

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type validatorService struct{}

type validatorServiceInterface interface {
	Get(name string) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	Create(dto *ValidatorDto) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	Update(*ValidatorDto, *ethereum2v1alpha1.Validator) (*ethereum2v1alpha1.Validator, *errors.RestErr)
	List() (*ethereum2v1alpha1.ValidatorList, *errors.RestErr)
	Delete(node *ethereum2v1alpha1.Validator) *errors.RestErr
	Count() (*int, *errors.RestErr)
}

var (
	ValidatorService validatorServiceInterface
)

func init() { ValidatorService = &validatorService{} }

// Get gets a single ethereum 2.0 beacon node by name
func (service validatorService) Get(name string) (*ethereum2v1alpha1.Validator, *errors.RestErr) {

	validator := &ethereum2v1alpha1.Validator{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, validator); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("validator by name %s doesn't exit", name))
		}
		go logger.Error("ERROR GETTING A VALIDATOR", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get a validator by name %s", name))
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
		},
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

	if err := k8s.Client().Create(context.Background(), validator); err != nil {
		if k8errors.IsAlreadyExists(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("validator by name %s already exits", validator.Name))
		}
		go logger.Error("ERROR_IN_CREATE_VALIDATOR", err)
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

	if err := k8s.Client().Update(context.Background(), validator); err != nil {
		go logger.Error("ERROR_IN_UPDATE_VALIDATOR", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update validator by name %s", validator.Name))
	}

	return validator, nil
}

// List returns all ethereum 2.0 beacon nodes
func (service validatorService) List() (*ethereum2v1alpha1.ValidatorList, *errors.RestErr) {
	validators := &ethereum2v1alpha1.ValidatorList{}

	if err := k8s.Client().List(context.Background(), validators, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_LIST_VALIDATOR", err)
		return nil, errors.NewInternalServerError("failed to get all validators")
	}

	return validators, nil
}

// Count returns total number of beacon nodes
func (service validatorService) Count() (*int, *errors.RestErr) {
	validators := &ethereum2v1alpha1.ValidatorList{}

	if err := k8s.Client().List(context.Background(), validators, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR COUNTING VALIDAITORS", err)
		return nil, errors.NewInternalServerError("error counting validators")
	}

	length := len(validators.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 beacon node by name
func (service validatorService) Delete(validator *ethereum2v1alpha1.Validator) *errors.RestErr {
	if err := k8s.Client().Delete(context.Background(), validator); err != nil {
		go logger.Error("ERROR_DELETING_VALIDATOR", err)
		return errors.NewBadRequestError(fmt.Sprintf("can't delete validator by name %s", validator.Name))
	}

	return nil
}
