package bitcoin

import (
	"context"
	"fmt"
	"github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	bitcointv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type bitcoinService struct{}

type IService interface {
	Get(types.NamespacedName) (*bitcointv1alpha1.Node, *errors.RestErr)
	List(namespace string) (*bitcointv1alpha1.NodeList, *errors.RestErr)
	Count(namespace string) (*int, *errors.RestErr)
	Delete(node *bitcointv1alpha1.Node) *errors.RestErr
}

var (
	k8sClient = k8s.NewClientService()
)

func NewBitcoinService() IService {
	return bitcoinService{}
}

// Get returns a single bitcoin node by name
func (service bitcoinService) Get(namespacedName types.NamespacedName) (*bitcointv1alpha1.Node, *errors.RestErr) {

	node := &bitcointv1alpha1.Node{}
	if err := k8sClient.Get(context.Background(), namespacedName, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", namespacedName.Name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", namespacedName.Name))
	}

	return node, nil
}

// List returns all bitcoin nodes
func (service bitcoinService) List(namespace string) (*bitcointv1alpha1.NodeList, *errors.RestErr) {
	nodes := &bitcointv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Count returns all nodes length
func (service bitcoinService) Count(namespace string) (*int, *errors.RestErr) {
	nodes := &bitcointv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	length := len(nodes.Items)
	return &length, nil
}

// Delete deletes bitcoin node by name
func (service bitcoinService) Delete(node *bitcointv1alpha1.Node) *errors.RestErr {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delte node by name %s", node.Name))
	}
	return nil
}
