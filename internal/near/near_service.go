package near

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type nearService struct{}

type nearServiceInterface interface {
	Get(name string) (*nearv1alpha1.Node, *restErrors.RestErr)
	Create(dto *NearDto) (*nearv1alpha1.Node, *restErrors.RestErr)
	Update(*NearDto, *nearv1alpha1.Node) (*nearv1alpha1.Node, *restErrors.RestErr)
	List() (*nearv1alpha1.NodeList, *restErrors.RestErr)
	Delete(node *nearv1alpha1.Node) *restErrors.RestErr
	Count() (*int, *restErrors.RestErr)
}

var (
	NearService nearServiceInterface
)

func init() { NearService = &nearService{} }

// Get gets a single filecoin node by name
func (service nearService) Get(name string) (*nearv1alpha1.Node, *restErrors.RestErr) {
	node := &nearv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("node by name %s doesn't exit", name))
		}
		go logger.Error(service.Get, err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil
}

// Create creates filecoin node from spec
func (service nearService) Create(dto *NearDto) (*nearv1alpha1.Node, *restErrors.RestErr) {
	node := &nearv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
		},
		Spec: nearv1alpha1.NodeSpec{
			Network: dto.Network,
			Archive: dto.Archive,
			RPC:     true,
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(context.Background(), node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("node by name %s already exits", node.Name))
		}
		go logger.Error(service.Create, err)
		return nil, restErrors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates filecoin node by name from spec
func (service nearService) Update(dto *NearDto, node *nearv1alpha1.Node) (*nearv1alpha1.Node, *restErrors.RestErr) {

	if dto.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = dto.NodePrivateKeySecretName
	}

	if dto.ValidatorSecretName != "" {
		node.Spec.ValidatorSecretName = dto.ValidatorSecretName
	}

	if dto.MinPeers != 0 {
		node.Spec.MinPeers = dto.MinPeers
	}

	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.P2PHost != "" {
		node.Spec.P2PHost = dto.P2PHost
	}

	if dto.RPC != nil {
		node.Spec.RPC = *dto.RPC
	}
	if node.Spec.RPC {
		if dto.RPCPort != 0 {
			node.Spec.RPCPort = dto.RPCPort
		}
		if dto.RPCHost != "" {
			node.Spec.RPCHost = dto.RPCHost
		}
	}

	if dto.PrometheusPort != 0 {
		node.Spec.PrometheusPort = dto.PrometheusPort
	}

	if dto.PrometheusHost != "" {
		node.Spec.PrometheusHost = dto.PrometheusHost
	}

	if dto.TelemetryURL != "" {
		node.Spec.TelemetryURL = dto.TelemetryURL
	}

	if bootnodes := dto.Bootnodes; bootnodes != nil {
		node.Spec.Bootnodes = *bootnodes
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
		go logger.Error(service.Update, err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all filecoin nodes
func (service nearService) List() (*nearv1alpha1.NodeList, *restErrors.RestErr) {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error(service.List, err)
		return nil, restErrors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns total number of filecoin nodes
func (service nearService) Count() (*int, *restErrors.RestErr) {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error(service.Count, err)
		return nil, restErrors.NewInternalServerError("failed to count all nodes")
	}

	length := len(nodes.Items)

	return &length, nil
}

// Delete deletes ethereum 2.0 filecoin node by name
func (service nearService) Delete(node *nearv1alpha1.Node) *restErrors.RestErr {
	if err := k8s.Client().Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
