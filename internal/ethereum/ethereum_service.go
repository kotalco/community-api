// Package ethereum internal is the domain layer for the ethereum node
// uses the k8 client to CRUD the node
package ethereum

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ethereumService struct{}

type ethereumServiceInterface interface {
	Get(name string) (*ethereumv1alpha1.Node, *errors.RestErr)
	Create(node *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr)
	Update(node *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr)
	List() (*ethereumv1alpha1.NodeList, *errors.RestErr)
	Delete(node *ethereumv1alpha1.Node) *errors.RestErr
}

var (
	EthereumService ethereumServiceInterface
)

func init() { EthereumService = &ethereumService{} }

// Get returns a single chainlink node by name
func (service ethereumService) Get(name string) (*ethereumv1alpha1.Node, *errors.RestErr) {
	node := &ethereumv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}
	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", name))
		}
		go logger.Error("Error Getting ethereum Node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil

}

// Create creates chainlink node from the given spec
func (service ethereumService) Create(node *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr) {
	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Create(context.Background(), node)
	if err != nil {
		if k8errors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
		}
		go logger.Error("error creating ethereum node", err)
		return nil, errors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates a single chainlink node by name from spec
func (service ethereumService) Update(node *ethereumv1alpha1.Node) (*ethereumv1alpha1.Node, *errors.RestErr) {
	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Update(context.Background(), node)
	if err != nil {
		go logger.Error("error updating ethereum node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all chainlink nodes
func (service ethereumService) List() (*ethereumv1alpha1.NodeList, *errors.RestErr) {
	nodes := &ethereumv1alpha1.NodeList{}

	err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default"))
	if err != nil {
		go logger.Error("Error listing ethereum nodes", err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Delete a single chainlink node by name
func (service ethereumService) Delete(node *ethereumv1alpha1.Node) *errors.RestErr {
	err := k8s.Client().Delete(context.Background(), node)

	if err != nil {
		go logger.Error("Error deleting ethereum node", err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
