package k8s

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/logger"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type MetaDataDto struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func (metaDto *MetaDataDto) ObjectMetaFromMetadataDto() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      metaDto.Name,
		Namespace: metaDto.Namespace,
	}
}

func DefaultResources(res *sharedAPI.Resources) {
	res.CPU = "1"
	res.Memory = "1Gi"
}

// CheckDeploymentResourcesChanged checks if any of the deployment resources passed by the users have changes in shareAPI.Resources
func CheckDeploymentResourcesChanged(res *sharedAPI.Resources) bool {
	if res.CPU != "" ||
		res.CPULimit != "" ||
		res.Memory != "" ||
		res.MemoryLimit != "" ||
		res.Storage != "" {
		return true
	}
	return false
}

// DeployReconciliation gets the pod related ot the deployment check if it still stuck in the pending state due to the insufficient resources,
// then we delete said pod, so it can get created again with updated (decreased) deployment specs resources
func DeployReconciliation(name string, namespace string) *restErrors.RestErr {
	var k8Client = NewClientService()
	// get pod
	pod := &corev1.Pod{}
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s-0", name),
	}
	err := k8Client.Get(context.Background(), key, pod)
	if apiErrors.IsNotFound(err) {
		go logger.Error("DEPLOY_RECONCILIATION", err)
		return restErrors.NewBadRequestError(fmt.Sprintf("pod by name %s doesn't exit", key.Name))
	}

	//check if pod status
	if pod.Status.Phase == "Pending" {
		err = k8Client.Delete(context.Background(), pod)
		if err != nil {
			go logger.Error("", err)
			return restErrors.NewInternalServerError(fmt.Sprintf("can't update deploy by name %s", name))
		}
	}
	return nil
}
