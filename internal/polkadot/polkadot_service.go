package polkadot

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type polkadtoService struct{}

type polkadotServiceInterface interface {
	Get(name string) (*polkadotv1alpha1.Node, *restErrors.RestErr)
	Create(dto *PolkadotDto) (*polkadotv1alpha1.Node, *restErrors.RestErr)
	Update(*PolkadotDto, *polkadotv1alpha1.Node) (*polkadotv1alpha1.Node, *restErrors.RestErr)
	List() (*polkadotv1alpha1.NodeList, *restErrors.RestErr)
	Delete(node *polkadotv1alpha1.Node) *restErrors.RestErr
	Count() (*int, *restErrors.RestErr)
}

var (
	PolkadotService polkadotServiceInterface
)

func init() { PolkadotService = &polkadtoService{} }

// Get gets a single filecoin node by name
func (service polkadtoService) Get(name string) (*polkadotv1alpha1.Node, *restErrors.RestErr) {
	node := &polkadotv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exits", name))
		}
		go logger.Error("ERROR_IN_GET_POLKADOT_SERVICE", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil
}

// Create creates filecoin node from spec
func (service polkadtoService) Create(dto *PolkadotDto) (*polkadotv1alpha1.Node, *restErrors.RestErr) {
	node := &polkadotv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
		},
		Spec: polkadotv1alpha1.NodeSpec{
			Network: dto.Network,
			RPC:     true,
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(context.Background(), node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("node by name %s is already exits", node.Name))
		}
		go logger.Error("ERROR_IN_CREATE_POLKADOT_SERVICE", err)
		return nil, restErrors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates filecoin node by name from spec
func (service polkadtoService) Update(dto *PolkadotDto, node *polkadotv1alpha1.Node) (*polkadotv1alpha1.Node, *restErrors.RestErr) {
	if dto.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = dto.NodePrivateKeySecretName
	}

	if dto.Validator != nil {
		node.Spec.Validator = *dto.Validator
	}

	if dto.SyncMode != "" {
		node.Spec.SyncMode = polkadotv1alpha1.SynchronizationMode(dto.SyncMode)
	}

	if dto.Pruning != nil {
		node.Spec.Pruning = dto.Pruning
	}

	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.RetainedBlocks != 0 {
		node.Spec.RetainedBlocks = dto.RetainedBlocks
	}

	if dto.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(dto.Logging)
	}

	if dto.Telemetry != nil {
		node.Spec.Telemetry = *dto.Telemetry
	}

	if dto.TelemetryURL != "" {
		node.Spec.TelemetryURL = dto.TelemetryURL
	}

	if dto.Prometheus != nil {
		node.Spec.Prometheus = *dto.Prometheus
	}

	if dto.PrometheusPort != 0 {
		node.Spec.PrometheusPort = dto.PrometheusPort
	}

	if dto.RPC != nil {
		node.Spec.RPC = *dto.RPC
	}

	if dto.RPCPort != 0 {
		node.Spec.RPCPort = dto.RPCPort
	}

	if dto.WS != nil {
		node.Spec.WS = *dto.WS
	}

	if dto.WSPort != 0 {
		node.Spec.WSPort = dto.WSPort
	}

	if len(dto.CORSDomains) != 0 {
		node.Spec.CORSDomains = dto.CORSDomains
	}

	if dto.CPU != "" {
		node.Spec.CPU = dto.CPU
	}

	if dto.CPULimit != "" {
		node.Spec.CPULimit = dto.CPULimit
	}

	if dto.Memory != "" {
		node.Spec.Memory = dto.Memory
	}

	if dto.MemoryLimit != "" {
		node.Spec.MemoryLimit = dto.MemoryLimit
	}

	if dto.Storage != "" {
		node.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Update(context.Background(), node); err != nil {
		go logger.Error("ERROR_IN_UPDATE_POLKADOT_SERVICE", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't updagte node by name %s", node.Name))
	}

	return node, nil
}

// List returns all filecoin nodes
func (service polkadtoService) List() (*polkadotv1alpha1.NodeList, *restErrors.RestErr) {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_LIST_POLKADOT_SERVICE", err)
		return nil, restErrors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns total number of filecoin nodes
func (service polkadtoService) Count() (*int, *restErrors.RestErr) {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_COUNT_SERVICE_POLKADOT", err)
		return nil, restErrors.NewInternalServerError("failed to count all nodes")
	}

	length := len(nodes.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 filecoin node by name
func (service polkadtoService) Delete(node *polkadotv1alpha1.Node) *restErrors.RestErr {
	if err := k8s.Client().Delete(context.Background(), node); err != nil {
		go logger.Error("ERROR_IN_DELETE_POLKADOT_SERVICE", err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delte node by name %s", node.Name))
	}

	return nil
}
