// Package storage_class internal is the domain layer for creating storage_class
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
	Get(name types.NamespacedName) (*storagev1.StorageClass, *errors.RestErr)
	Create(dto *StorageClassDto) (*storagev1.StorageClass, *errors.RestErr)
	Update(*StorageClassDto, *storagev1.StorageClass) (*storagev1.StorageClass, *errors.RestErr)
	List(namespace string) (*storagev1.StorageClassList, *errors.RestErr)
	Delete(*storagev1.StorageClass) *errors.RestErr
	Count(namespace string) (*int, *errors.RestErr)
}

var (
	StorageClassService storageClassServiceInterface
	k8Client            = k8s.K8ClientService
)

func init() { StorageClassService = &storageClassService{} }

// Get returns a single storage class  by name
func (service storageClassService) Get(namespacedName types.NamespacedName) (*storagev1.StorageClass, *errors.RestErr) {
	storageClass := &storagev1.StorageClass{}

	if err := k8Client.Get(context.Background(), namespacedName, storageClass); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("storage class by name %s doens't exit", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get storage class by name %s", namespacedName.Name))
	}

	return storageClass, nil
}

// Create creates a storage class from the given spec
//todo
func (service storageClassService) Create(dto *StorageClassDto) (*storagev1.StorageClass, *errors.RestErr) {
	return nil, nil
}

// Update creates a storage class from the given spec
//todo
func (service storageClassService) Update(dto *StorageClassDto, storageClass *storagev1.StorageClass) (*storagev1.StorageClass, *errors.RestErr) {
	return nil, nil
}

// List returns all storage classes
func (service storageClassService) List(namespace string) (*storagev1.StorageClassList, *errors.RestErr) {
	storageClasses := &storagev1.StorageClassList{}

	if err := k8Client.List(context.Background(), storageClasses, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get storage class list")
	}

	return storageClasses, nil
}

// Delete a single storage node by name
//todo
func (service storageClassService) Delete(storageClass *storagev1.StorageClass) *errors.RestErr {
	return nil
}

// Count a list of storage classes
//todo
func (service storageClassService) Count(namespace string) (*int, *errors.RestErr) {
	return nil, nil
}
