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
)

type bitcoinService struct{}

type IService interface {
	Get(types.NamespacedName) (*bitcointv1alpha1.Node, *errors.RestErr)
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
