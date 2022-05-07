package storage_class

import (
	"github.com/kotalco/api/pkg/k8s"
	storagev1 "k8s.io/api/storage/v1"
)

// StorageClass is Kubernetes storage class
type StorageClassDto struct {
	k8s.MetaDataDto
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
}

type StorageClassListDto []StorageClassDto

func (dto StorageClassDto) FromCoreStorageClass(sc *storagev1.StorageClass) *StorageClassDto {
	dto.Name = sc.Name
	dto.Provisioner = sc.Provisioner
	dto.ReclaimPolicy = string(*sc.ReclaimPolicy)
	dto.AllowVolumeExpansion = sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion

	return &dto
}

func (storageClassListDto StorageClassListDto) FromCoreSecret(list []storagev1.StorageClass) StorageClassListDto {
	result := make(StorageClassListDto, len(list))
	for index, value := range list {
		result[index] = *(StorageClassDto{}.FromCoreStorageClass(&value))
	}
	return result
}
