// Package storage_class internal is the domain layer for creating storage_class
// uses the k8 client to CRUD the storage_class
package storage_class

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	storagev1 "k8s.io/api/storage/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type storageClassService struct{}

type IService interface {
	Get(name string) (storagev1.StorageClass, restErrors.IRestErr)
	Create(dto StorageClassDto) (storagev1.StorageClass, restErrors.IRestErr)
	Update(StorageClassDto, *storagev1.StorageClass) restErrors.IRestErr
	List() (storagev1.StorageClassList, restErrors.IRestErr)
	Delete(*storagev1.StorageClass) restErrors.IRestErr
	Count() (int, restErrors.IRestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewStorageClassService() IService {
	return storageClassService{}
}

// Get returns a single storage class  by name
func (service storageClassService) Get(name string) (storagev1.StorageClass, restErrors.IRestErr) {
	storageClass := &storagev1.StorageClass{}
	key := types.NamespacedName{Name: name}
	if err := k8sClient.Get(context.Background(), key, storageClass); err != nil {
		if apiErrors.IsNotFound(err) {
			return storagev1.StorageClass{}, restErrors.NewNotFoundError(fmt.Sprintf("storage class by name %s doens't exit", key.Name))
		}
		go logger.Error(service.Get, err)
		return storagev1.StorageClass{}, restErrors.NewInternalServerError(fmt.Sprintf("can't get storage class by name %s", key.Name))
	}

	return *storageClass, nil
}

// Create creates a storage class from the given spec
// todo
func (service storageClassService) Create(dto StorageClassDto) (storagev1.StorageClass, restErrors.IRestErr) {
	return storagev1.StorageClass{}, nil
}

// Update creates a storage class from the given spec
// todo
func (service storageClassService) Update(dto StorageClassDto, storageClass *storagev1.StorageClass) restErrors.IRestErr {
	return nil
}

// List returns all storage classes
func (service storageClassService) List() (storagev1.StorageClassList, restErrors.IRestErr) {
	storageClasses := &storagev1.StorageClassList{}

	if err := k8sClient.List(context.Background(), storageClasses); err != nil {
		go logger.Error(service.List, err)
		return storagev1.StorageClassList{}, restErrors.NewInternalServerError("failed to get storage class list")
	}

	return *storageClasses, nil
}

// Delete a single storage node by name
// todo
func (service storageClassService) Delete(storageClass *storagev1.StorageClass) restErrors.IRestErr {
	return nil
}

// Count a list of storage classes
// todo
func (service storageClassService) Count() (int, restErrors.IRestErr) {
	return 0, nil
}
