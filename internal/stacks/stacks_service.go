package stacks

import (
	"context"
	"fmt"
	"github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	stacksv1alpha1 "github.com/kotalco/kotal/apis/stacks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type stacksService struct{}

type IService interface {
	Get(types.NamespacedName) (*stacksv1alpha1.Node, *errors.RestErr)
	List(namespace string) (*stacksv1alpha1.NodeList, *errors.RestErr)
	Count(namespace string) (*int, *errors.RestErr)
	Delete(node *stacksv1alpha1.Node) *errors.RestErr
	Update(dto *StacksDto, node *stacksv1alpha1.Node) (*stacksv1alpha1.Node, *errors.RestErr)
}

var (
	k8sClient = k8s.NewClientService()
)

func NewStacksService() IService {
	return stacksService{}
}

// Get returns a single stacks node by name
func (service stacksService) Get(namespacedName types.NamespacedName) (*stacksv1alpha1.Node, *errors.RestErr) {

	node := &stacksv1alpha1.Node{}
	if err := k8sClient.Get(context.Background(), namespacedName, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
	}

	return node, nil
}

// List returns all stacks nodes
func (service stacksService) List(namespace string) (*stacksv1alpha1.NodeList, *errors.RestErr) {
	nodes := &stacksv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns all nodes length
func (service stacksService) Count(namespace string) (*int, *errors.RestErr) {
	nodes := &stacksv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	length := len(nodes.Items)
	return &length, nil
}

// Update updates a single node by name from spec
func (service stacksService) Update(dto *StacksDto, node *stacksv1alpha1.Node) (*stacksv1alpha1.Node, *errors.RestErr) {
	if dto.Image != "" {
		node.Spec.Image = dto.Image
	}
	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}
	if dto.RPCPort != 0 {
		node.Spec.RPCPort = dto.RPCPort
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
			return nil, errors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
		}
		podIsPending = pod.Status.Phase == corev1.PodPending
	}

	if err := k8sClient.Update(context.Background(), node); err != nil {
		go logger.Error(service.Update, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	if podIsPending {
		err := k8sClient.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error(service.Update, err)
			return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
		}
	}

	return node, nil
}

// Delete deletes stacks node by name
func (service stacksService) Delete(node *stacksv1alpha1.Node) *errors.RestErr {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}
	return nil
}
