// Package secret internal is the domain layer for creating storage_class
// uses the k8 client to CRUD the storage_class
package storage_class

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	storagev1 "k8s.io/api/storage/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type storageClassService struct{}

type storageClassServiceInterface interface {
	Get(name string) (*storagev1.StorageClass, *errors.RestErr)
	Create(dto *StorageClassDto) (*storagev1.StorageClass, *errors.RestErr)
	Update(*StorageClassDto, *storagev1.StorageClass) (*storagev1.StorageClass, *errors.RestErr)
	List() (*storagev1.StorageClassList, *errors.RestErr)
	Delete(*storagev1.StorageClass) *errors.RestErr
	Count() (*int, *errors.RestErr)
}

var (
	StorageClassService storageClassServiceInterface
)

func init() { StorageClassService = &storageClassService{} }

// Get returns a single secret  by name
func (service storageClassService) Get(name string) (*storagev1.StorageClass, *errors.RestErr) {
	storageClass := &storagev1.StorageClass{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, storageClass); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("storage class by name %s doens't exit", name))
		}
		go logger.Error("ERROR_IN_GET_STORAGE_CLASS", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get storage class by name %s", name))
	}

	return storageClass, nil
}

// Create creates a secret from the given spec
func (service storageClassService) Create(dto *StorageClassDto) (*storagev1.StorageClass, *errors.RestErr) {
	return nil, nil
}

// Create creates a secret from the given spec
func (service storageClassService) Update(dto *StorageClassDto, storageClass *storagev1.StorageClass) (*storagev1.StorageClass, *errors.RestErr) {
	return nil, nil
}

// List returns all secrets
func (service storageClassService) List() (*storagev1.StorageClassList, *errors.RestErr) {
	storageClasses := &storagev1.StorageClassList{}

	if err := k8s.Client().List(context.Background(), storageClasses, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_LIST_STORAGE_CLASS", err)
		return nil, errors.NewInternalServerError("failed to get storage class list")
	}

	return storageClasses, nil
}

// Delete a single secret node by name
func (service storageClassService) Delete(storageClass *storagev1.StorageClass) *errors.RestErr {
	return nil
}

// Delete a list of secrets
func (service storageClassService) Count() (*int, *errors.RestErr) {
	return nil, nil
}