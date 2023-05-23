package bitcoin

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	bitcoinv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type bitcoinService struct{}

type IService interface {
	Get(types.NamespacedName) (bitcoinv1alpha1.Node, restErrors.IRestErr)
	List(namespace string) (bitcoinv1alpha1.NodeList, restErrors.IRestErr)
	Count(namespace string) (int, restErrors.IRestErr)
	Create(BitcoinDto) (bitcoinv1alpha1.Node, restErrors.IRestErr)
	Delete(*bitcoinv1alpha1.Node) restErrors.IRestErr
	Update(BitcoinDto, *bitcoinv1alpha1.Node) restErrors.IRestErr
}

var (
	k8sClient = k8s.NewClientService()
)

func NewBitcoinService() IService {
	return bitcoinService{}
}

// Get returns a single bitcoin node by name
func (service bitcoinService) Get(namespacedName types.NamespacedName) (node bitcoinv1alpha1.Node, restErr restErrors.IRestErr) {
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

// List returns all bitcoin nodes
func (service bitcoinService) List(namespace string) (list bitcoinv1alpha1.NodeList, restErr restErrors.IRestErr) {
	err := k8sClient.List(context.Background(), &list, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.List, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}
	return
}

// Count returns all nodes length
func (service bitcoinService) Count(namespace string) (count int, restErr restErrors.IRestErr) {
	nodes := &bitcoinv1alpha1.NodeList{}
	err := k8sClient.List(context.Background(), nodes, client.InNamespace(namespace))
	if err != nil {
		go logger.Error(service.Count, err)
		restErr = restErrors.NewInternalServerError("failed to get all nodes")
		return
	}

	return len(nodes.Items), nil
}

// Create creates bitcoin node from the given specs
func (service bitcoinService) Create(dto BitcoinDto) (node bitcoinv1alpha1.Node, restErr restErrors.IRestErr) {
	node.ObjectMeta = dto.ObjectMetaFromMetadataDto()
	node.Spec = bitcoinv1alpha1.NodeSpec{
		Network: dto.Network,
		RPC:     true,
		RPCUsers: []bitcoinv1alpha1.RPCUser{
			{
				Username:           BitcoinJsonRpcDefaultUserName,
				PasswordSecretName: BitcoinJsonRpcDefaultUserPasswordName,
			},
		},
	}

	k8s.DefaultResources(&node.Spec.Resources)

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

// Update updates a single node by name from spec
func (service bitcoinService) Update(dto BitcoinDto, node *bitcoinv1alpha1.Node) (restErr restErrors.IRestErr) {
	if dto.Image != "" {
		node.Spec.Image = dto.Image
	}
	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}
	if dto.RPC != nil {
		node.Spec.RPC = *dto.RPC
	}
	if len(dto.RPCUsers) != 0 {
		node.Spec.RPCUsers = make([]bitcoinv1alpha1.RPCUser, 0)
		for _, v := range dto.RPCUsers {
			node.Spec.RPCUsers = append(node.Spec.RPCUsers, bitcoinv1alpha1.RPCUser{
				Username:           v.Username,
				PasswordSecretName: v.PasswordSecretName,
			})
		}
	}
	if dto.Wallet != nil {
		node.Spec.Wallet = *dto.Wallet
	}
	if dto.TransactionIndex != nil {
		node.Spec.TransactionIndex = *dto.TransactionIndex
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

// Delete deletes bitcoin node by name
func (service bitcoinService) Delete(node *bitcoinv1alpha1.Node) (restErr restErrors.IRestErr) {
	if err := k8sClient.Delete(context.Background(), node); err != nil {
		go logger.Error(service.Delete, err)
		restErr = restErrors.NewInternalServerError(fmt.Sprintf("can't delete node by name %s", node.Name))
		return
	}
	return
}
