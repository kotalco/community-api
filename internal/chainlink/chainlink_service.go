// Package chainlink internal is the domain layer for the chainlink node
// uses the k8 client to CRUD the node
package chainlink

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type chainlinkService struct{}

type chainlinkServiceInterface interface {
	Get(name string) (*chainlinkv1alpha1.Node, *errors.RestErr)
	Create(node *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr)
	Update(node *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr)
	List() (*chainlinkv1alpha1.NodeList, *errors.RestErr)
	Delete(node *chainlinkv1alpha1.Node) *errors.RestErr
}

var (
	ChainlinkService chainlinkServiceInterface
)

func init() { ChainlinkService = &chainlinkService{} }

// Get returns a single chainlink node by name
func (service chainlinkService) Get(name string) (*chainlinkv1alpha1.Node, *errors.RestErr) {
	node := &chainlinkv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}
	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exist", name))
		}
		go logger.Error("Error Getting Node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil

}

// Create creates chainlink node from the given spec
func (service chainlinkService) Create(node *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr) {
	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Create(context.Background(), node)
	if err != nil {
		if k8errors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("node by name %s already exist", node.Name))
		}
		go logger.Error("error creating chainlink node", err)
		return nil, errors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates a single chainlink node by name from spec
func (service chainlinkService) Update(node *chainlinkv1alpha1.Node) (*chainlinkv1alpha1.Node, *errors.RestErr) {
	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	err := k8s.Client().Update(context.Background(), node)
	if err != nil {
		go logger.Error("error updating chainlink node", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all chainlink nodes
func (service chainlinkService) List() (*chainlinkv1alpha1.NodeList, *errors.RestErr) {
	nodes := &chainlinkv1alpha1.NodeList{}

	err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default"))
	if err != nil {
		go logger.Error("Error listing chainlink nodes", err)
		return nil, errors.NewInternalServerError("failed to get all nodes")
	}

	return nodes, nil
}

// Delete a single chainlink node by name
func (service chainlinkService) Delete(node *chainlinkv1alpha1.Node) *errors.RestErr {
	err := k8s.Client().Delete(context.Background(), node)

	if err != nil {
		go logger.Error("Error deleting chainlink node", err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
	}

	return nil
}
