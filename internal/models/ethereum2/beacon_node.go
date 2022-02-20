package models

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

type BeaconNode struct {
	models.Time
	Name    string `json:"name"`
	Network string `json:"network"`
	Client  string `json:"client"`
	// todo: required only for prysm and network is not mainnet
	Eth1Endpoints *[]string `json:"eth1Endpoints"`
	REST          *bool     `json:"rest"`
	RESTHost      string    `json:"restHost"`
	RESTPort      uint      `json:"restPort"`
	RPC           *bool     `json:"rpc"`
	RPCHost       string    `json:"rpcHost"`
	RPCPort       uint      `json:"rpcPort"`
	GRPC          *bool     `json:"grpc"`
	GRPCHost      string    `json:"grpcHost"`
	GRPCPort      uint      `json:"grpcPort"`
	CPU           string    `json:"cpu"`
	CPULimit      string    `json:"cpuLimit"`
	Memory        string    `json:"memory"`
	MemoryLimit   string    `json:"memoryLimit"`
	Storage       string    `json:"storage"`
	StorageClass  *string   `json:"storageClass"`
}

func FromEthereum2BeaconNode(beaconnode *ethereum2v1alpha1.BeaconNode) *BeaconNode {
	return &BeaconNode{
		Name: beaconnode.Name,
		Time: models.Time{
			CreatedAt: beaconnode.CreationTimestamp.UTC().Format(shared.JavascriptISOString),
		},
		Network:       beaconnode.Spec.Network,
		Client:        string(beaconnode.Spec.Client),
		Eth1Endpoints: &beaconnode.Spec.Eth1Endpoints,
		REST:          &beaconnode.Spec.REST,
		RESTHost:      beaconnode.Spec.RESTHost,
		RESTPort:      beaconnode.Spec.RESTPort,
		RPC:           &beaconnode.Spec.RPC,
		RPCHost:       beaconnode.Spec.RPCHost,
		RPCPort:       beaconnode.Spec.RPCPort,
		GRPC:          &beaconnode.Spec.GRPC,
		GRPCHost:      beaconnode.Spec.GRPCHost,
		GRPCPort:      beaconnode.Spec.GRPCPort,
		CPU:           beaconnode.Spec.CPU,
		CPULimit:      beaconnode.Spec.CPULimit,
		Memory:        beaconnode.Spec.Memory,
		MemoryLimit:   beaconnode.Spec.MemoryLimit,
		Storage:       beaconnode.Spec.Storage,
		StorageClass:  beaconnode.Spec.StorageClass,
	}
}
