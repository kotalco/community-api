package chainlink

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

type apiCredentials struct {
	Email              string `json:"email"`
	PasswordSecretName string `json:"passwordSecretName"`
}

type ChainlinkDto struct {
	models.Time
	k8s.MetaDataDto
	EthereumChainId            uint            `json:"ethereumChainId"`
	LinkContractAddress        string          `json:"linkContractAddress"`
	EthereumWSEndpoint         string          `json:"ethereumWsEndpoint"`
	DatabaseURL                string          `json:"databaseURL"`
	EthereumHTTPEndpoints      []string        `json:"ethereumHttpEndpoints"`
	KeystorePasswordSecretName string          `json:"keystorePasswordSecretName"`
	APICredentials             *apiCredentials `json:"apiCredentials"`
	CORSDomains                []string        `json:"corsDomains"`
	CertSecretName             string          `json:"certSecretName"`
	TLSPort                    uint            `json:"tlsPort"`
	P2PPort                    uint            `json:"p2pPort"`
	APIPort                    uint            `json:"apiPort"`
	SecureCookies              *bool           `json:"secureCookies"`
	Logging                    string          `json:"logging"`
	sharedAPI.Resources
}

type ChainlinkListDto []ChainlinkDto

func (dto ChainlinkDto) FromChainlinkNode(n *chainlinkv1alpha1.Node) *ChainlinkDto {
	dto.Name = n.Name
	dto.Time = models.Time{CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.EthereumChainId = n.Spec.EthereumChainId
	dto.LinkContractAddress = n.Spec.LinkContractAddress
	dto.EthereumWSEndpoint = n.Spec.EthereumWSEndpoint
	dto.DatabaseURL = n.Spec.DatabaseURL
	dto.EthereumHTTPEndpoints = n.Spec.EthereumHTTPEndpoints
	dto.KeystorePasswordSecretName = n.Spec.KeystorePasswordSecretName
	dto.APICredentials = &apiCredentials{
		Email:              n.Spec.APICredentials.Email,
		PasswordSecretName: n.Spec.APICredentials.PasswordSecretName,
	}
	dto.CORSDomains = n.Spec.CORSDomains
	dto.CertSecretName = n.Spec.CertSecretName
	dto.TLSPort = n.Spec.TLSPort
	dto.P2PPort = n.Spec.P2PPort
	dto.APIPort = n.Spec.APIPort
	dto.SecureCookies = &n.Spec.SecureCookies
	dto.Logging = string(n.Spec.Logging)
	dto.CPU = n.Spec.CPU
	dto.CPULimit = n.Spec.CPULimit
	dto.Memory = n.Spec.Memory
	dto.MemoryLimit = n.Spec.MemoryLimit
	dto.Storage = n.Spec.Storage
	dto.StorageClass = n.Spec.StorageClass

	return &dto
}

func (nodes ChainlinkListDto) FromChainlinkNode(models []chainlinkv1alpha1.Node) ChainlinkListDto {
	result := make(ChainlinkListDto, len(models))
	for index, model := range models {
		result[index] = *(ChainlinkDto{}.FromChainlinkNode(&model))
	}
	return result
}
