package statefulset

import (
	"context"
	"fmt"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

var k8sClient = k8s.NewClientService()

type IStatefulSet interface {
	Get(namespacedName types.NamespacedName) (*appsv1.StatefulSet, restError.IRestErr)
}

type statefulset struct {
}

func NewService() IStatefulSet {
	return &statefulset{}
}

func (s *statefulset) Get(namespacedName types.NamespacedName) (*appsv1.StatefulSet, restError.IRestErr) {
	record := &appsv1.StatefulSet{}

	err := k8sClient.Get(context.Background(), namespacedName, record)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, restError.NewNotFoundError(fmt.Sprintf("record with the name %s doesn't exist", namespacedName.Name))
		}
		go logger.Error(s.Get, err)
		return nil, restError.NewInternalServerError("can't list stateful set")
	}
	return record, nil
}
