package aptos

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	aptosv1alpha1 "github.com/kotalco/kotal/apis/aptos/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type aptosService struct{}

type IService interface {
	// Get returns a single aptos node by name
	Get(types.NamespacedName) (aptosv1alpha1.Node, restErrors.IRestErr)
	// List returns all aptos nodes
	List(namespace string) (aptosv1alpha1.NodeList, restErrors.IRestErr)
	// Count returns all nodes length
	Count(namespace string) (int, restErrors.IRestErr)
	// Create creates aptos node from the given specs
	Create(AptosDto) (aptosv1alpha1.Node, restErrors.IRestErr)
	// Delete deletes aptos node by name
	Delete(*aptosv1alpha1.Node) restErrors.IRestErr
	// Update updates a single node by name from spec
	Update(AptosDto, *aptosv1alpha1.Node) restErrors.IRestErr
}

var (
	k8sClient = k8s.NewClientService()
)

func NewAptosService() IService {
	return aptosService{}
}

func (service aptosService) Get(namespacedName types.NamespacedName) (node aptosv1alpha1.Node, restErr restErrors.IRestErr) {
	if err := k8sClient.Get(context.Background(), namespacedName, &node); err != nil {
		if apiErrors.IsNotFound(err) {
			restErr = restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", namespacedName.Name))
			return
		}
		go logger.Error(service.Get, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
		return
	}
	return
}

func (service aptosService) List(namespace string) (list aptosv1alpha1.NodeList, restErr restErrors.IRestErr) {
	err := k8sClient.List(context.Background(), &list, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}
	return
}

func (service aptosService) Count(namespace string) (count int, restErr restErrors.IRestErr) {
	nodes := &aptosv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}

	return len(nodes.Items), nil
}

func (service aptosService) Create(dto AptosDto) (node aptosv1alpha1.Node, restErr restErrors.IRestErr) {
	node.ObjectMeta = dto.ObjectMetaFromMetadataDto()
	k8s.DefaultResources(&node.Spec.Resources)
	node.Spec.Network = dto.Network
	node.Spec.Image = dto.Image
	node.Spec.API = true

	if err := k8sClient.Create(context.Background(), &node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			restErr = restErrors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
			return
		}
		go logger.Error(service.Create, err)
		restErr = restErrors.NewInternalServerError("failed to create node")
		return
	}

	return
}

func (service aptosService) Update(dto AptosDto, node *aptosv1alpha1.Node) (restErr restErrors.IRestErr) {
	if dto.Image != "" {
		node.Spec.Image = dto.Image
	}
	if dto.Validator != nil {
		node.Spec.Validator = *dto.Validator
	}

	if dto.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = dto.NodePrivateKeySecretName
	}

	if dto.API != nil {
		node.Spec.API = *dto.API
	}
	if dto.APIPort != 0 {
		node.Spec.APIPort = dto.APIPort
	}
	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
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
			restErr = restErrors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
			return
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
		return
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
			return
		}
	}

	return
}

func (service aptosService) Delete(node *aptosv1alpha1.Node) (restErr restErrors.IRestErr) {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
		return
	}
	return
}
