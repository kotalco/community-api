package models

import (
	storagev1 "k8s.io/api/storage/v1"
)

// StorageClass is Kubernetes storage class
type StorageClass struct {
	Name                 string `json:"name"`
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
}

func FromCoreStorageClass(sc *storagev1.StorageClass) *StorageClass {
	return &StorageClass{
		Name:                 sc.Name,
		Provisioner:          sc.Provisioner,
		ReclaimPolicy:        string(*sc.ReclaimPolicy),
		AllowVolumeExpansion: sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion,
	}
}
