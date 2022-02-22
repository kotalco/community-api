package chainlink

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
)

type apiCredentials struct {
	Email              string `json:"email"`
	PasswordSecretName string `json:"passwordSecretName"`
}

type ChainlinkDto struct {
	models.Time
	Name                       string          `json:"name"`
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
	CPU                        string          `json:"cpu"`
	CPULimit                   string          `json:"cpuLimit"`
	Memory                     string          `json:"memory"`
	MemoryLimit                string          `json:"memoryLimit"`
	Storage                    string          `json:"storage"`
	StorageClass               *string         `json:"storageClass"`
}

type ChainlinkListDto []ChainlinkDto

func (node ChainlinkDto) FromChainlinkNode(n *chainlinkv1alpha1.Node) *ChainlinkDto {
	return &ChainlinkDto{
		Name: n.Name,
		Time: models.Time{
			CreatedAt: n.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		EthereumChainId:            n.Spec.EthereumChainId,
		LinkContractAddress:        n.Spec.LinkContractAddress,
		EthereumWSEndpoint:         n.Spec.EthereumWSEndpoint,
		DatabaseURL:                n.Spec.DatabaseURL,
		EthereumHTTPEndpoints:      n.Spec.EthereumHTTPEndpoints,
		KeystorePasswordSecretName: n.Spec.KeystorePasswordSecretName,
		APICredentials: &apiCredentials{
			Email:              n.Spec.APICredentials.Email,
			PasswordSecretName: n.Spec.APICredentials.PasswordSecretName,
		},
		CORSDomains:    n.Spec.CORSDomains,
		CertSecretName: n.Spec.CertSecretName,
		TLSPort:        n.Spec.TLSPort,
		P2PPort:        n.Spec.P2PPort,
		APIPort:        n.Spec.APIPort,
		SecureCookies:  &n.Spec.SecureCookies,
		Logging:        string(n.Spec.Logging),
		CPU:            n.Spec.CPU,
		CPULimit:       n.Spec.CPULimit,
		Memory:         n.Spec.Memory,
		MemoryLimit:    n.Spec.MemoryLimit,
		Storage:        n.Spec.Storage,
		StorageClass:   n.Spec.StorageClass,
	}
}

func (nodes ChainlinkListDto) FromChainlinkNode(models []chainlinkv1alpha1.Node) ChainlinkListDto {
	result := make(ChainlinkListDto, len(models))
	for index, model := range models {
		result[index] = *(ChainlinkDto{}.FromChainlinkNode(&model))
	}
	return result
}
