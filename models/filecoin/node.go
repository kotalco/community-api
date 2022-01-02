package models

import (
	"github.com/kotalco/api/models"
	"github.com/kotalco/api/shared"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
)

// Node is Filecoin node
type Node struct {
	models.Time
	Name         string  `json:"name"`
	Network      string  `json:"network"`
	CPU          string  `json:"cpu"`
	CPULimit     string  `json:"cpuLimit"`
	Memory       string  `json:"memory"`
	MemoryLimit  string  `json:"memoryLimit"`
	Storage      string  `json:"storage"`
	StorageClass *string `json:"storageClass"`
}

// FromFilecoinNode creates node model from Filecoin node
func FromFilecoinNode(node *filecoinv1alpha1.Node) *Node {
	return &Node{
		Name: node.Name,
		Time: models.Time{
			CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:      string(node.Spec.Network),
		CPU:          node.Spec.CPU,
		CPULimit:     node.Spec.CPULimit,
		Memory:       node.Spec.Memory,
		MemoryLimit:  node.Spec.MemoryLimit,
		Storage:      node.Spec.Storage,
		StorageClass: node.Spec.StorageClass,
	}
}
