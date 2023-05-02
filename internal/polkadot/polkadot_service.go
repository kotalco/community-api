package polkadot

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type polkadtoService struct{}

type IService interface {
	Get(types.NamespacedName) (polkadotv1alpha1.Node, *restErrors.RestErr)
	Create(PolkadotDto) (polkadotv1alpha1.Node, *restErrors.RestErr)
	Update(PolkadotDto, *polkadotv1alpha1.Node) *restErrors.RestErr
	List(namespace string) (polkadotv1alpha1.NodeList, *restErrors.RestErr)
	Delete(*polkadotv1alpha1.Node) *restErrors.RestErr
	Count(namespace string) (int, *restErrors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewPolkadotService() IService {
	return polkadtoService{}
}

// Get gets a single filecoin node by name
func (service polkadtoService) Get(namespacedName types.NamespacedName) (polkadotv1alpha1.Node, *restErrors.RestErr) {
	node := &polkadotv1alpha1.Node{}

	if err := k8sClient.Get(context.Background(), namespacedName, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return polkadotv1alpha1.Node{}, restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exits", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return polkadotv1alpha1.Node{}, restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
	}

	return *node, nil
}

// Create creates filecoin node from spec
func (service polkadtoService) Create(dto PolkadotDto) (polkadotv1alpha1.Node, *restErrors.RestErr) {
	node := &polkadotv1alpha1.Node{
		ObjectMeta: dto.ObjectMetaFromMetadataDto(),
		Spec: polkadotv1alpha1.NodeSpec{
			Network: dto.Network,
			RPC:     true,
			Pruning: dto.Pruning,
			Image:   dto.Image,
		},
	}

	k8s.DefaultResources(&node.Spec.Resources)

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8sClient.Create(context.Background(), node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return polkadotv1alpha1.Node{}, restErrors.NewBadRequestError(fmt.Sprintf("node by name %s is already exits", node.Name))
		}
		go logger.Error(service.Create, err)
		return polkadotv1alpha1.Node{}, restErrors.NewInternalServerError("failed to create node")
	}

	return *node, nil
}

// Update updates filecoin node by name from spec
func (service polkadtoService) Update(dto PolkadotDto, node *polkadotv1alpha1.Node) *restErrors.RestErr {
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
	if dto.Image != "" {
		node.Spec.Image = dto.Image
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	pod := &corev1.Pod{}
	podIsPending := false
	if dto.CPU != "" || dto.Memory != "" {
		key := types.NamespacedName{
			Namespace: node.Namespace,
			Name:      fmt.Sprintf("%s-0", node.Name),
		}
		err := k8sClient.Get(context.Background(), key, pod)
		if apiErrors.IsNotFound(err) {
			go logger.Error(service.Update, err)
			return restErrors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			return restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
		}
	}
	return nil
}

// List returns all filecoin nodes
func (service polkadtoService) List(namespace string) (polkadotv1alpha1.NodeList, *restErrors.RestErr) {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.List, err)
		return polkadotv1alpha1.NodeList{}, restErrors.NewInternalServerError("failed to get all nodes")
	}

	return *nodes, nil
}

// Count returns total number of filecoin nodes
func (service polkadtoService) Count(namespace string) (int, *restErrors.RestErr) {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace)); err != nil {
		go logger.Error(service.Count, err)
		return 0, restErrors.NewInternalServerError("failed to count all nodes")
	}

	return len(nodes.Items), nil
}

// Delete deletes ethereum 2.0 filecoin node by name
func (service polkadtoService) Delete(node *polkadotv1alpha1.Node) *restErrors.RestErr {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delte node by name %s", node.Name))
	}

	return nil
}
