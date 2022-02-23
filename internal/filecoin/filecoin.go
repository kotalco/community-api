package filecoin

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
)

// Node is Filecoin node
type FilecoinDto struct {
	models.Time
	Name               string  `json:"name"`
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
func (filecoinDto FilecoinDto) FromFilecoinNode(node *filecoinv1alpha1.Node) *FilecoinDto {
	return &FilecoinDto{
		Name: node.Name,
		Time: models.Time{
			CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:            string(node.Spec.Network),
		API:                &node.Spec.API,
		APIPort:            node.Spec.APIPort,
		APIHost:            node.Spec.APIHost,
		APIRequestTimeout:  node.Spec.APIRequestTimeout,
		DisableMetadataLog: &node.Spec.DisableMetadataLog,
		P2PPort:            node.Spec.P2PPort,
		P2PHost:            node.Spec.P2PHost,
		IPFSPeerEndpoint:   node.Spec.IPFSPeerEndpoint,
		IPFSOnlineMode:     &node.Spec.IPFSOnlineMode,
		IPFSForRetrieval:   &node.Spec.IPFSForRetrieval,
		CPU:                node.Spec.CPU,
		CPULimit:           node.Spec.CPULimit,
		Memory:             node.Spec.Memory,
		MemoryLimit:        node.Spec.MemoryLimit,
		Storage:            node.Spec.Storage,
		StorageClass:       node.Spec.StorageClass,
	}
}

// FromFilecoinNode creates node dto from Filecoin node list
func (filecoinListDto FilecoinListDto) FromFilecoinNode(nodes []filecoinv1alpha1.Node) FilecoinListDto {
	result := make(FilecoinListDto, len(nodes))
	for index, v := range nodes {
		result[index] = *(FilecoinDto{}.FromFilecoinNode(&v))
	}
	return result
}