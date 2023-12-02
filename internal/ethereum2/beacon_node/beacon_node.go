package beacon_node

import (
	"github.com/kotalco/community-api/internal/models"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
)

type BeaconNodeDto struct {
	models.Time
	k8s.MetaDataDto
	Network                 string  `json:"network"`
	Client                  string  `json:"client"`
	REST                    *bool   `json:"rest"`
	RESTPort                uint    `json:"restPort"`
	RPC                     *bool   `json:"rpc"`
	RPCPort                 uint    `json:"rpcPort"`
	GRPC                    *bool   `json:"grpc"`
	GRPCPort                uint    `json:"grpcPort"`
	ExecutionEngineEndpoint *string `json:"executionEngineEndpoint"`
	CheckpointSyncURL       *string `json:"checkpointSyncUrl"`
	JWTSecretName           *string `json:"jwtSecretName"`
	Image                   string  `json:"image"`
	sharedAPI.Resources
}
type BeaconNodeListDto []BeaconNodeDto

func (dto BeaconNodeDto) FromEthereum2BeaconNode(node ethereum2v1alpha1.BeaconNode) BeaconNodeDto {
	dto.Name = node.Name
	dto.Time = models.Time{CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = node.Spec.Network
	dto.Client = string(node.Spec.Client)
	dto.REST = &node.Spec.REST
	dto.RESTPort = node.Spec.RESTPort
	dto.RPC = &node.Spec.RPC
	dto.RPCPort = node.Spec.RPCPort
	dto.GRPC = &node.Spec.GRPC
	dto.GRPCPort = node.Spec.GRPCPort
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass
	dto.ExecutionEngineEndpoint = &node.Spec.ExecutionEngineEndpoint
	dto.JWTSecretName = &node.Spec.JWTSecretName
	dto.Image = node.Spec.Image
	dto.CheckpointSyncURL = &node.Spec.CheckpointSyncURL

	return dto
}

func (nodes BeaconNodeListDto) FromEthereum2BeaconNode(beaconnodeList []ethereum2v1alpha1.BeaconNode) BeaconNodeListDto {
	result := make(BeaconNodeListDto, len(beaconnodeList))
	for index, v := range beaconnodeList {
		result[index] = BeaconNodeDto{}.FromEthereum2BeaconNode(v)
	}
	return result
}
