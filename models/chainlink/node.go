package models

import chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"

type APICredentials struct {
	Email              string `json:"email"`
	PasswordSecretName string `json:"passwordSecretName"`
}

type Node struct {
	Name                       string          `json:"name"`
	EthereumChainId            uint            `json:"ethereumChainId"`
	LinkContractAddress        string          `json:"linkContractAddress"`
	EthereumWSEndpoint         string          `json:"ethereumWsEndpoint"`
	DatabaseURL                string          `json:"databaseURL"`
	EthereumHTTPEndpoints      []string        `json:"ethereumHttpEndpoints"`
	KeystorePasswordSecretName string          `json:"keystorePasswordSecretName"`
	APICredentials             *APICredentials `json:"apiCredentials"`
	CORSDomains                []string        `json:"corsDomains"`
	CertSecretName             string          `json:"certSecretName"`
	TLSPort                    uint            `json:"tlsPort"`
	P2PPort                    uint            `json:"p2pPort"`
	APIPort                    uint            `json:"apiPort"`
}

func FromChainlinkNode(node *chainlinkv1alpha1.Node) *Node {
	return &Node{
		Name:                       node.Name,
		EthereumChainId:            node.Spec.EthereumChainId,
		LinkContractAddress:        node.Spec.LinkContractAddress,
		EthereumWSEndpoint:         node.Spec.EthereumWSEndpoint,
		DatabaseURL:                node.Spec.DatabaseURL,
		EthereumHTTPEndpoints:      node.Spec.EthereumHTTPEndpoints,
		KeystorePasswordSecretName: node.Spec.KeystorePasswordSecretName,
		APICredentials: &APICredentials{
			Email:              node.Spec.APICredentials.Email,
			PasswordSecretName: node.Spec.APICredentials.PasswordSecretName,
		},
		CORSDomains:    node.Spec.CORSDomains,
		CertSecretName: node.Spec.CertSecretName,
		TLSPort:        node.Spec.TLSPort,
		P2PPort:        node.Spec.P2PPort,
		APIPort:        node.Spec.APIPort,
	}
}
