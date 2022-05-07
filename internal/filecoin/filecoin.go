package filecoin

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
)

// Node is Filecoin node
type FilecoinDto struct {
	models.Time
	k8s.MetaDataDto
	Network            string  `json:"network"`
	API                *bool   `json:"api"`
	APIPort            uint    `json:"apiPort"`
	APIHost            string  `json:"apiHost"`
	APIRequestTimeout  uint    `json:"apiRequestTimeout"`
	DisableMetadataLog *bool   `json:"disableMetadataLog"`
	P2PPort            uint    `json:"p2pPort"`
	P2PHost            string  `json:"p2pHost"`
	IPFSPeerEndpoint   string  `json:"ipfsPeerEndpoint"`
	IPFSOnlineMode     *bool   `json:"ipfsOnlineMode"`
	IPFSForRetrieval   *bool   `json:"ipfsForRetrieval"`
	CPU                string  `json:"cpu"`
	CPULimit           string  `json:"cpuLimit"`
	Memory             string  `json:"memory"`
	MemoryLimit        string  `json:"memoryLimit"`
	Storage            string  `json:"storage"`
	StorageClass       *string `json:"storageClass"`
}

type FilecoinListDto []FilecoinDto

// FromFilecoinNode creates node dto from Filecoin node
func (dto FilecoinDto) FromFilecoinNode(node *filecoinv1alpha1.Node) *FilecoinDto {

	dto.Name = node.Name
	dto.Time = models.Time{CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = string(node.Spec.Network)
	dto.API = &node.Spec.API
	dto.APIPort = node.Spec.APIPort
	dto.APIHost = node.Spec.APIHost
	dto.APIRequestTimeout = node.Spec.APIRequestTimeout
	dto.DisableMetadataLog = &node.Spec.DisableMetadataLog
	dto.P2PPort = node.Spec.P2PPort
	dto.P2PHost = node.Spec.P2PHost
	dto.IPFSPeerEndpoint = node.Spec.IPFSPeerEndpoint
	dto.IPFSOnlineMode = &node.Spec.IPFSOnlineMode
	dto.IPFSForRetrieval = &node.Spec.IPFSForRetrieval
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass

	return &dto
}

// FromFilecoinNode creates node dto from Filecoin node list
func (filecoinListDto FilecoinListDto) FromFilecoinNode(nodes []filecoinv1alpha1.Node) FilecoinListDto {
	result := make(FilecoinListDto, len(nodes))
	for index, v := range nodes {
		result[index] = *(FilecoinDto{}.FromFilecoinNode(&v))
	}
	return result
}
