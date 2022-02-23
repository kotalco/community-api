package storage_class

import (
	storagev1 "k8s.io/api/storage/v1"
)

// StorageClass is Kubernetes storage class
type StorageClassDto struct {
	Name                 string `json:"name"`
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
}

type StorageClassListDto []StorageClassDto

func (storageClass StorageClassDto) FromCoreStorageClass(sc *storagev1.StorageClass) *StorageClassDto {
	return &StorageClassDto{
		Name:                 sc.Name,
		Provisioner:          sc.Provisioner,
		ReclaimPolicy:        string(*sc.ReclaimPolicy),
		AllowVolumeExpansion: sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion,
	}
}

func (storageClassListDto StorageClassListDto) FromCoreSecret(list []storagev1.StorageClass) StorageClassListDto {
	result := make(StorageClassListDto, len(list))
	for index, value := range list {
		result[index] = *(StorageClassDto{}.FromCoreStorageClass(&value))
	}
	return result
}
