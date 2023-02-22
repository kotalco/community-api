package shared

import (
	"context"
	"errors"
	"fmt"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/shared"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type metricsResponseDto struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

// Metrics returns a websocket that emits cpu and memory usage
func Metrics(c *websocket.Conn) {
	defer c.Close()

	pod := &corev1.Pod{}
	key := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      fmt.Sprintf("%s-0", c.Params("name")),
	}

	sts := &appsv1.StatefulSet{}
	stsKey := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      c.Params("name"),
	}

	for {

		err := k8sClient.Get(context.Background(), key, pod)
		if err != nil {
			// is the pod error due to sts has been deleted ?
			stsErr := k8sClient.Get(context.Background(), stsKey, sts)
			if apierrors.IsNotFound(stsErr) {
				c.WriteJSON(shared.NewResponse(restError.NewNotFoundError(stsErr.Error())))
				return
			}

			if apierrors.IsNotFound(err) {
				err = errors.New("NotFound")
			}
			c.WriteJSON(shared.NewResponse(restError.NewNotFoundError(err.Error())))

			time.Sleep(3 * time.Second)
			continue
		}

		metricsClientset := k8s.MetricsClientset()
		opts := metav1.GetOptions{}
		metrics, err := metricsClientset.MetricsV1beta1().PodMetricses(key.Namespace).Get(context.Background(), key.Name, opts)
		if err != nil {
			c.WriteJSON(restError.NewInternalServerError(err.Error()))
			return
		}

		response := new(metricsResponseDto)
		//cpu are represented in nano-cores which is 1/1000000000 (1 billionth) of a cpu
		//scaling to milli core which is 1/1000 of a cpu
		response.Cpu = metrics.Containers[0].Usage.Cpu().ScaledValue(resource.Milli)
		//memory usage are represented in ki  (1 Kibibyte = 1.024 kilobytes) (1000 Kibibyte  = 1.024 megabytes)
		//scaling to megabytes
		//the value won't overflow int64 coz we are scaling to megabytes
		response.Memory = metrics.Containers[0].Usage.Memory().ScaledValue(resource.Mega)

		c.WriteJSON(shared.NewResponse(response))
		time.Sleep(time.Second * 3)
	}

}
