package statefulset

import (
	"context"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	k8sClient = k8s.NewClientService()
)

type IStatefulSet interface {
	Exists(name string) (bool, *restError.RestErr)
}

type stateful struct {
}

func NewService() IStatefulSet {
	return &stateful{}
}

func (s *stateful) Exists(name string) (bool, *restError.RestErr) {
	list := &appsv1.StatefulSetList{}

	err := k8sClient.List(context.Background(), list, &client.MatchingLabels{"app.kubernetes.io/managed-by": "kotal-operator"})
	if err != nil {
		go logger.Error(s.Exists, err)
		return false, restError.NewInternalServerError("can't list stateful set")
	}
	for _, v := range list.Items {
		if v.Name == name {
			return true, nil
		}
	}
	return false, nil
}
