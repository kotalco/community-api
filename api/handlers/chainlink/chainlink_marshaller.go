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

type Node struct {
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

type Nodes []Node

func (node Node) FromChainlinkNode(model *chainlinkv1alpha1.Node) *Node {
	return &Node{
		Name: model.Name,
		Time: models.Time{
			CreatedAt: model.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		EthereumChainId:            model.Spec.EthereumChainId,
		LinkContractAddress:        model.Spec.LinkContractAddress,
		EthereumWSEndpoint:         model.Spec.EthereumWSEndpoint,
		DatabaseURL:                model.Spec.DatabaseURL,
		EthereumHTTPEndpoints:      model.Spec.EthereumHTTPEndpoints,
		KeystorePasswordSecretName: model.Spec.KeystorePasswordSecretName,
		APICredentials: &apiCredentials{
			Email:              model.Spec.APICredentials.Email,
			PasswordSecretName: model.Spec.APICredentials.PasswordSecretName,
		},
		CORSDomains:    model.Spec.CORSDomains,
		CertSecretName: model.Spec.CertSecretName,
		TLSPort:        model.Spec.TLSPort,
		P2PPort:        model.Spec.P2PPort,
		APIPort:        model.Spec.APIPort,
		SecureCookies:  &model.Spec.SecureCookies,
		Logging:        string(model.Spec.Logging),
		CPU:            model.Spec.CPU,
		CPULimit:       model.Spec.CPULimit,
		Memory:         model.Spec.Memory,
		MemoryLimit:    model.Spec.MemoryLimit,
		Storage:        model.Spec.Storage,
		StorageClass:   model.Spec.StorageClass,
	}
}

func (nodes Nodes) FromChainlinkNode(models []chainlinkv1alpha1.Node) Nodes {
	result := make(Nodes, len(models))
	for index, model := range models {
		result[index] = *(Node{}.FromChainlinkNode(&model))
	}
	return result
}
